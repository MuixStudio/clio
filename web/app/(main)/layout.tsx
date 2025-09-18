"use client";
import React from "react";

import Header from "@/components/header";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen flex-col">
      <div className="shrink-0">
        <Header />
      </div>
      <div className="flex h-full w-full flex-1 overflow-y-auto">
        {children}
      </div>
    </div>
  );
}
