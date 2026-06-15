import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Field, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import {
  Popover,
  PopoverContent,
  PopoverDescription,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"

import { PaletteGroup, PaletteSection } from "./palette-section"

export function OverlaysSection() {
  return (
    <PaletteSection
      id="overlays"
      title="弹层 Overlays"
      description="dialog.tsx, alert-dialog.tsx, sheet.tsx, popover.tsx, dropdown-menu.tsx"
    >
      <PaletteGroup label="dialog">
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="outline">打开 Dialog</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>编辑资料</DialogTitle>
              <DialogDescription>修改完成后点击保存</DialogDescription>
            </DialogHeader>
            <Field>
              <FieldLabel htmlFor="dialog-name">名称</FieldLabel>
              <Input id="dialog-name" defaultValue="Pedro Duarte" />
            </Field>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">取消</Button>
              </DialogClose>
              <Button>保存</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </PaletteGroup>

      <PaletteGroup label="alert dialog">
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive">删除账户</Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>确定要删除吗？</AlertDialogTitle>
              <AlertDialogDescription>
                此操作无法撤销，将永久删除该账户的所有数据。
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>取消</AlertDialogCancel>
              <AlertDialogAction>确定删除</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </PaletteGroup>

      <PaletteGroup label="sheet">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="outline">打开 Sheet</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>编辑资料</SheetTitle>
              <SheetDescription>在侧边栏中修改信息，完成后点击保存</SheetDescription>
            </SheetHeader>
            <Field className="px-4">
              <FieldLabel htmlFor="sheet-name">名称</FieldLabel>
              <Input id="sheet-name" defaultValue="Pedro Duarte" />
            </Field>
            <SheetFooter>
              <Button>保存</Button>
              <SheetClose asChild>
                <Button variant="outline">取消</Button>
              </SheetClose>
            </SheetFooter>
          </SheetContent>
        </Sheet>
      </PaletteGroup>

      <PaletteGroup label="popover">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline">打开 Popover</Button>
          </PopoverTrigger>
          <PopoverContent>
            <PopoverHeader>
              <PopoverTitle>尺寸设置</PopoverTitle>
              <PopoverDescription>调整组件的宽高</PopoverDescription>
            </PopoverHeader>
            <div className="grid grid-cols-2 gap-2">
              <Field>
                <FieldLabel htmlFor="popover-width">宽度</FieldLabel>
                <Input id="popover-width" defaultValue="100%" />
              </Field>
              <Field>
                <FieldLabel htmlFor="popover-height">高度</FieldLabel>
                <Input id="popover-height" defaultValue="auto" />
              </Field>
            </div>
          </PopoverContent>
        </Popover>
      </PaletteGroup>

      <PaletteGroup label="dropdown menu">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">打开菜单</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>我的账户</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              个人信息
              <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuItem>账单</DropdownMenuItem>
            <DropdownMenuCheckboxItem checked>
              显示状态栏
            </DropdownMenuCheckboxItem>
            <DropdownMenuSeparator />
            <DropdownMenuRadioGroup value="light">
              <DropdownMenuRadioItem value="light">浅色主题</DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="dark">深色主题</DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem variant="destructive">退出登录</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </PaletteGroup>
    </PaletteSection>
  )
}
