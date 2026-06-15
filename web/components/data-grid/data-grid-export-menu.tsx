"use client";

import type { Row, Table } from "@tanstack/react-table";
import { Download, FileJson2, FileSpreadsheet, FileText } from "lucide-react";
import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface DataGridExportMenuProps<TData> {
  table: Table<TData>;
  filename?: string;
  disabled?: boolean;
}

export function DataGridExportMenu<TData>({
  table,
  filename = "export",
  disabled,
}: DataGridExportMenuProps<TData>) {
  const selectedCount = table.getSelectedRowModel().rows.length;

  const getVisibleColumns = React.useCallback(
    () =>
      table
        .getAllColumns()
        .filter(
          (col) =>
            col.getIsVisible() &&
            typeof col.accessorFn !== "undefined" &&
            col.id !== "select",
        ),
    [table],
  );

  const toMatrix = React.useCallback(
    (rows: Row<TData>[]) => {
      const cols = getVisibleColumns();
      const headers = cols.map((col) => col.columnDef.meta?.label ?? col.id);
      const data = rows.map((row) =>
        cols.map((col) => {
          const cell = row.getAllCells().find((c) => c.column.id === col.id);
          const value = cell?.getValue();
          if (Array.isArray(value)) return value.join(", ");
          if (value === null || value === undefined) return "";
          return value as string | number | boolean;
        }),
      );
      return { headers, data };
    },
    [getVisibleColumns],
  );

  const toObjects = React.useCallback(
    (rows: Row<TData>[]) => {
      const cols = getVisibleColumns();
      return rows.map((row) => {
        const obj: Record<string, unknown> = {};
        for (const col of cols) {
          const key = col.columnDef.meta?.label ?? col.id;
          const cell = row.getAllCells().find((c) => c.column.id === col.id);
          obj[key] = cell?.getValue() ?? null;
        }
        return obj;
      });
    },
    [getVisibleColumns],
  );

  const exportExcel = React.useCallback(
    async (rows: Row<TData>[], suffix = "") => {
      const { headers, data } = toMatrix(rows);
      const XLSX = await import("xlsx");
      const ws = XLSX.utils.aoa_to_sheet([headers, ...data]);
      const wb = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, "Sheet1");
      XLSX.writeFile(wb, `${filename}${suffix}.xlsx`);
    },
    [toMatrix, filename],
  );

  const exportCsv = React.useCallback(
    (rows: Row<TData>[], suffix = "") => {
      const { headers, data } = toMatrix(rows);
      const lines = [
        headers.map(escapeCsv).join(","),
        ...data.map((row) => row.map((cell) => escapeCsv(String(cell))).join(",")),
      ];
      triggerDownload(
        new Blob([lines.join("\n")], { type: "text/csv;charset=utf-8;" }),
        `${filename}${suffix}.csv`,
      );
    },
    [toMatrix, filename],
  );

  const exportJson = React.useCallback(
    (rows: Row<TData>[], suffix = "") => {
      triggerDownload(
        new Blob([JSON.stringify(toObjects(rows), null, 2)], {
          type: "application/json",
        }),
        `${filename}${suffix}.json`,
      );
    },
    [toObjects, filename],
  );

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="font-bold"
          disabled={disabled}
        >
          <Download className="text-muted-foreground" />
          Export
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel className="text-xs font-normal text-muted-foreground">
          All rows
        </DropdownMenuLabel>
        <DropdownMenuGroup>
          <DropdownMenuItem
            onSelect={() => exportExcel(table.getFilteredRowModel().rows)}
          >
            <FileSpreadsheet />
            Excel
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={() => exportCsv(table.getFilteredRowModel().rows)}
          >
            <FileText />
            CSV
          </DropdownMenuItem>
          <DropdownMenuItem
            onSelect={() => exportJson(table.getFilteredRowModel().rows)}
          >
            <FileJson2 />
            JSON
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuLabel className="text-xs font-normal text-muted-foreground">
          {selectedCount > 0
            ? `Selected rows (${selectedCount})`
            : "Selected rows"}
        </DropdownMenuLabel>
        <DropdownMenuGroup>
          <DropdownMenuItem
            disabled={selectedCount === 0}
            onSelect={() =>
              exportExcel(table.getSelectedRowModel().rows, "_selected")
            }
          >
            <FileSpreadsheet />
            Excel
          </DropdownMenuItem>
          <DropdownMenuItem
            disabled={selectedCount === 0}
            onSelect={() =>
              exportCsv(table.getSelectedRowModel().rows, "_selected")
            }
          >
            <FileText />
            CSV
          </DropdownMenuItem>
          <DropdownMenuItem
            disabled={selectedCount === 0}
            onSelect={() =>
              exportJson(table.getSelectedRowModel().rows, "_selected")
            }
          >
            <FileJson2 />
            JSON
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function escapeCsv(value: string): string {
  if (/[",\n\r]/.test(value)) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

function triggerDownload(blob: Blob, name: string): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  a.click();
  URL.revokeObjectURL(url);
}