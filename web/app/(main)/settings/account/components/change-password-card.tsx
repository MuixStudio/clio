import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

export function ChangePasswordCard() {
  return (
    <Card className="w-full">
      <CardHeader className="">
        <CardTitle>Change Password</CardTitle>
        <CardDescription>change you password</CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="checkout-7j9-card-name-43j">
                Set Password
              </FieldLabel>
              <FieldDescription>
                Enter your 16-digit password here
              </FieldDescription>
              <Input
                id="checkout-7j9-card-name-43j"
                placeholder="********"
                required
              />
              <FieldLabel htmlFor="checkout-7j9-card-name-43j">
                Confirm password
              </FieldLabel>
              <Input
                id="checkout-7j9-card-name-43j"
                placeholder="********"
                required
              />
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
      <CardFooter className="flex-raw gap-2">
        <Button type="submit" className="">
          Change Password
        </Button>
      </CardFooter>
    </Card>
  )
}
