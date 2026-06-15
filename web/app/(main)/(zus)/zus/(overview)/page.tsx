"use client"

import { ChartAreaInteractive } from "./charts/chart"
import {
  HeaderBar,
  HeaderBarDescription,
  HeaderBarLeft,
  HeaderBarTitle,
} from "@/app/(main)/components/header-bar/header-bar"

export default function Page() {
  return (
    <div>
      {/*header*/}
      <HeaderBar>
        <HeaderBarLeft>
          <HeaderBarTitle className="flex items-center gap-2">
            Overview
          </HeaderBarTitle>
          <HeaderBarDescription>Overview</HeaderBarDescription>
        </HeaderBarLeft>
      </HeaderBar>

      <br />

      <ChartAreaInteractive></ChartAreaInteractive>
    </div>
  )
}
