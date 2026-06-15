import { CreditCard, Search } from "lucide-react"

import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"

import { PaletteGroup, PaletteSection } from "./palette-section"

export function InputsSection() {
  return (
    <PaletteSection
      id="inputs"
      title="输入 Input"
      description="input.tsx, textarea.tsx, input-group.tsx, field.tsx"
    >
      <PaletteGroup label="input">
        <Input placeholder="Default" className="w-48" />
        <Input placeholder="Disabled" disabled className="w-48" />
        <Input placeholder="Invalid" aria-invalid className="w-48" />
        <Input type="password" placeholder="Password" className="w-48" />
      </PaletteGroup>

      <PaletteGroup label="textarea">
        <Textarea placeholder="Type your message..." className="w-72" />
      </PaletteGroup>

      <PaletteGroup label="input group" className="items-stretch">
        <InputGroup className="w-72">
          <InputGroupInput placeholder="Search..." />
          <InputGroupAddon>
            <Search />
          </InputGroupAddon>
        </InputGroup>
        <InputGroup className="w-72">
          <InputGroupAddon>
            <CreditCard />
          </InputGroupAddon>
          <InputGroupInput placeholder="Card number" />
          <InputGroupAddon align="inline-end">
            <InputGroupButton size="xs">验证</InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </PaletteGroup>

      <PaletteGroup label="field">
        <FieldSet className="w-full max-w-md">
          <FieldLegend>个人信息</FieldLegend>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="palette-name">姓名</FieldLabel>
              <Input id="palette-name" placeholder="请输入姓名" />
              <FieldDescription>显示在个人主页上的名称</FieldDescription>
            </Field>
            <Field data-invalid="true">
              <FieldLabel htmlFor="palette-email">邮箱</FieldLabel>
              <Input
                id="palette-email"
                type="email"
                aria-invalid
                placeholder="you@example.com"
              />
              <FieldError errors={[{ message: "请输入有效的邮箱地址" }]} />
            </Field>
          </FieldGroup>
        </FieldSet>
      </PaletteGroup>

      <PaletteGroup label="label">
        <Label htmlFor="palette-label-demo">Email address</Label>
        <Input id="palette-label-demo" placeholder="you@example.com" className="w-48" />
      </PaletteGroup>
    </PaletteSection>
  )
}
