import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Slider } from "@/components/ui/slider"

import { PaletteGroup, PaletteSection } from "./palette-section"

export function SelectionSection() {
  return (
    <PaletteSection
      id="selection"
      title="选择 Selection"
      description="select.tsx, checkbox.tsx, radio-group.tsx, slider.tsx"
    >
      <PaletteGroup label="select">
        <Select defaultValue="banana">
          <SelectTrigger className="w-44">
            <SelectValue placeholder="Default size" />
          </SelectTrigger>
          <SelectContent position="popper">
            <SelectGroup>
              <SelectLabel>Fruits</SelectLabel>
              <SelectItem value="apple">Apple</SelectItem>
              <SelectItem value="banana">Banana</SelectItem>
              <SelectItem value="orange">Orange</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select defaultValue="banana">
          <SelectTrigger size="sm" className="w-44">
            <SelectValue placeholder="Small size" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Fruits</SelectLabel>
              <SelectItem value="apple">Apple</SelectItem>
              <SelectItem value="banana">Banana</SelectItem>
              <SelectItem value="orange">Orange</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select disabled>
          <SelectTrigger className="w-44">
            <SelectValue placeholder="Disabled" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="apple">Apple</SelectItem>
          </SelectContent>
        </Select>
      </PaletteGroup>

      <PaletteGroup label="checkbox" className="flex-col items-start gap-2.5">
        <div className="flex items-center gap-2">
          <Checkbox id="palette-checkbox-default" />
          <Label htmlFor="palette-checkbox-default">Default</Label>
        </div>
        <div className="flex items-center gap-2">
          <Checkbox id="palette-checkbox-checked" defaultChecked />
          <Label htmlFor="palette-checkbox-checked">Checked</Label>
        </div>
        <div className="flex items-center gap-2">
          <Checkbox id="palette-checkbox-disabled" disabled />
          <Label htmlFor="palette-checkbox-disabled">Disabled</Label>
        </div>
        <div className="flex items-center gap-2">
          <Checkbox id="palette-checkbox-invalid" aria-invalid />
          <Label htmlFor="palette-checkbox-invalid">Invalid</Label>
        </div>
      </PaletteGroup>

      <PaletteGroup label="radio group" className="flex-col items-start">
        <RadioGroup defaultValue="comfortable" className="gap-2.5">
          <div className="flex items-center gap-2">
            <RadioGroupItem value="default" id="palette-radio-default" />
            <Label htmlFor="palette-radio-default">Default</Label>
          </div>
          <div className="flex items-center gap-2">
            <RadioGroupItem
              value="comfortable"
              id="palette-radio-comfortable"
            />
            <Label htmlFor="palette-radio-comfortable">Comfortable</Label>
          </div>
          <div className="flex items-center gap-2">
            <RadioGroupItem
              value="compact"
              id="palette-radio-compact"
              disabled
            />
            <Label htmlFor="palette-radio-compact">Compact (disabled)</Label>
          </div>
        </RadioGroup>
      </PaletteGroup>

      <PaletteGroup label="slider" className="flex-col items-stretch">
        <Slider defaultValue={[40]} className="w-64" />
        <Slider defaultValue={[20, 70]} className="w-64" />
        <Slider defaultValue={[40]} disabled className="w-64" />
      </PaletteGroup>
    </PaletteSection>
  )
}
