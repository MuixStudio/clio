"use client"

import { Button } from "@/components/ui/button"
import React from "react"

export default function Page() {
  const [n, setN] = React.useState(0)
  return (
    <div>
      <Button
        onClick={() => {
          setN(n+1)
        }}
      >
        aa
      </Button>
      {n}
    </div>
  )
}
