import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Field, FieldDescription, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

export function ChangeUserNameCard() {
  return (
    <Card className="w-full">
      <CardHeader className="">
        <CardTitle>Change UserName</CardTitle>
        <CardDescription>change you username</CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="checkout-7j9-card-name-43j">
                Name on Card
              </FieldLabel>
              <FieldDescription>
                Enter your 16-digit card number
              </FieldDescription>
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
      <CardFooter className="flex-raw gap-2">
        <Button type="submit" className="">
          Change Email
        </Button>
      </CardFooter>
    </Card>
  )
}