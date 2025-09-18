import React from "react";

type SettingCardGroupProps = {
  children: React.ReactNode;
};

export const SettingCardGroup = ({ children }: SettingCardGroupProps) => {
  return (
    <>
      <div className="w-full flex flex-col gap-6">{children}</div>
    </>
  );
};
