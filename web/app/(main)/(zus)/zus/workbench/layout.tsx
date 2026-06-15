import { WorkbenchSidebar } from "@/components/workbench-sidebar/workbench-sidebar"
import {
  SecondarySidebarInset,
  SecondarySidebarProvider,
} from "@/components/ui/secondary-sidebar"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SecondarySidebarProvider>
      <WorkbenchSidebar />
      <SecondarySidebarInset>{children}</SecondarySidebarInset>
    </SecondarySidebarProvider>
  )
}
