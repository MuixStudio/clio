import {
  Avatar,
  AvatarFallback,
  AvatarGroup,
  AvatarGroupCount,
  AvatarImage,
  AvatarBadge,
} from "@/components/ui/avatar"
import { Kbd, KbdGroup } from "@/components/ui/kbd"
import { Separator } from "@/components/ui/separator"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"

import { PaletteGroup, PaletteSection } from "./palette-section"

const INVOICES = [
  { invoice: "INV001", status: "已支付", method: "信用卡", amount: "￥250.00" },
  { invoice: "INV002", status: "待支付", method: "支付宝", amount: "￥150.00" },
  { invoice: "INV003", status: "未支付", method: "银行转账", amount: "￥350.00" },
]

export function DataDisplaySection() {
  return (
    <PaletteSection
      id="data-display"
      title="数据展示 Data Display"
      description="avatar.tsx, table.tsx, skeleton.tsx, kbd.tsx, separator.tsx, tooltip.tsx"
    >
      <PaletteGroup label="avatar — size">
        <Avatar size="sm">
          <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
          <AvatarFallback>SM</AvatarFallback>
        </Avatar>
        <Avatar size="default">
          <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
          <AvatarFallback>DF</AvatarFallback>
        </Avatar>
        <Avatar size="lg">
          <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
          <AvatarFallback>LG</AvatarFallback>
        </Avatar>
        <Avatar size="default">
          <AvatarFallback>状</AvatarFallback>
          <AvatarBadge className="bg-success" />
        </Avatar>
      </PaletteGroup>

      <PaletteGroup label="avatar group">
        <AvatarGroup>
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
            <AvatarFallback>A</AvatarFallback>
          </Avatar>
          <Avatar>
            <AvatarFallback>B</AvatarFallback>
          </Avatar>
          <Avatar>
            <AvatarFallback>C</AvatarFallback>
          </Avatar>
          <AvatarGroupCount>+5</AvatarGroupCount>
        </AvatarGroup>
      </PaletteGroup>

      <PaletteGroup label="kbd">
        <Kbd>Esc</Kbd>
        <KbdGroup>
          <Kbd>Ctrl</Kbd>
          <Kbd>K</Kbd>
        </KbdGroup>
        <Tooltip>
          <TooltipTrigger asChild>
            <span className="text-sm text-muted-foreground">悬浮查看快捷键</span>
          </TooltipTrigger>
          <TooltipContent side="top">
            <KbdGroup>
              <Kbd>Ctrl</Kbd>
              <Kbd>S</Kbd>
            </KbdGroup>
          </TooltipContent>
        </Tooltip>
      </PaletteGroup>

      <PaletteGroup label="separator" className="flex-col items-stretch">
        <div className="text-sm text-muted-foreground">水平分割线</div>
        <Separator />
        <div className="flex h-8 items-center gap-3 text-sm text-muted-foreground">
          <span>左</span>
          <Separator orientation="vertical" />
          <span>右</span>
        </div>
      </PaletteGroup>

      <PaletteGroup label="skeleton" className="flex-col items-start">
        <div className="flex items-center gap-3">
          <Skeleton className="size-10 rounded-full" />
          <div className="flex flex-col gap-2">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-3 w-24" />
          </div>
        </div>
      </PaletteGroup>

      <PaletteGroup label="table" className="flex-col items-stretch">
        <Table>
          <TableCaption>最近的账单记录</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead className="w-28">发票</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>支付方式</TableHead>
              <TableHead className="text-right">金额</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {INVOICES.map((invoice) => (
              <TableRow key={invoice.invoice}>
                <TableCell className="font-medium">{invoice.invoice}</TableCell>
                <TableCell>{invoice.status}</TableCell>
                <TableCell>{invoice.method}</TableCell>
                <TableCell className="text-right">{invoice.amount}</TableCell>
              </TableRow>
            ))}
          </TableBody>
          <TableFooter>
            <TableRow>
              <TableCell colSpan={3}>合计</TableCell>
              <TableCell className="text-right">￥750.00</TableCell>
            </TableRow>
          </TableFooter>
        </Table>
      </PaletteGroup>
    </PaletteSection>
  )
}
