import { Geist, Geist_Mono, Roboto } from "next/font/google"

import "@/styles/globals.css"
import { Providers } from "./providers"
import { cn } from "@/lib/utils"

const roboto = Roboto({ subsets: ["latin"], variable: "--font-sans" })

const fontMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
})

export const metadata = {
  title: "Clio",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={cn(
        "antialiased",
        fontMono.variable,
        "font-sans",
        roboto.variable
      )}
    >
      <body className={cn("overflow-hidden")}>
        <Providers>
          <main className="max-h-screen max-w-full">{children}</main>
        </Providers>
      </body>
    </html>
  )
}
