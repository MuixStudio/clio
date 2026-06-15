import { Toolbar, ToolbarContent, ToolbarItem } from "@/components/ui/toolbar"
import { DataGridFilterMenu } from "@/components/data-grid/data-grid-filter-menu"
import { DataGridSortMenu } from "@/components/data-grid/data-grid-sort-menu"
import { DataGridRowHeightMenu } from "@/components/data-grid/data-grid-row-height-menu"
import { DataGridViewMenu } from "@/components/data-grid/data-grid-view-menu"
import { DataGridExportMenu } from "@/components/data-grid/data-grid-export-menu"
import * as React from "react"
import { Table } from "@tanstack/react-table"

export default function DataGridToolbar<TData>({ table }: { table: Table<TData> }) {
  return (
    <Toolbar>
      <ToolbarContent>
        <ToolbarItem>
          <DataGridFilterMenu table={table} />
        </ToolbarItem>
        <ToolbarItem>
          <DataGridSortMenu table={table} />
        </ToolbarItem>
        <ToolbarItem>
          <DataGridRowHeightMenu table={table} />
        </ToolbarItem>
      </ToolbarContent>

      <ToolbarContent justify="end">
        <ToolbarItem>
          <DataGridViewMenu table={table} />
        </ToolbarItem>
        <ToolbarItem>
          <DataGridExportMenu table={table} />
        </ToolbarItem>
      </ToolbarContent>
    </Toolbar>
  )
}