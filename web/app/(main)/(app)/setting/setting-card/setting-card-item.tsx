import React from "react";

type SettingItemCardProps = {
  title?: string;
  description?: string;
  children: React.ReactNode;
};

export const SettingCardItem = ({
  title,
  description,
  children,
}: SettingItemCardProps) => {
  return (
    <>
      <div className="flex w-full flex-col gap-0.5">
        <p className="text-small font-medium">{title}</p>
        <p className="text-foreground-500 text-xs">{description}</p>
        <div>{children}</div>
      </div>
    </>
  );
};
