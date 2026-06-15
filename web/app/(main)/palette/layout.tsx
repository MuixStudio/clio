"use client"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto max-w-362 p-6">
      <div className="w-full">{children}</div>
    </div>
  )
}
