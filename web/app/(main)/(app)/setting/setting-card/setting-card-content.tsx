import React from "react";

type SettingContentProps = {
  children: React.ReactNode;
};

export const SettingCardContent = ({ children }: SettingContentProps) => {
  return (
    <>
      <div className="flex w-full flex-col gap-6">{children}</div>
    </>
  );
};
