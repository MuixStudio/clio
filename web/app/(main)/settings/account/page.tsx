import { ChangeEmailCard } from "./components/change-email-card"
import { ChangePasswordCard } from "./components/change-password-card"
import { ChangeUserNameCard } from "@/app/(main)/settings/account/components/change-username-card"

export default function Page() {
  return (
    <div className="flex w-full flex-col gap-4">
      <ChangeEmailCard />
      <ChangePasswordCard />
      <ChangeUserNameCard />
    </div>
  )
}
