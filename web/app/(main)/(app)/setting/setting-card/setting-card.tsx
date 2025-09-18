import React from "react";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";

type SettingItemProps = {
  title: string;
  description?: string;
  children: React.ReactNode;
};

export const SettingCard = ({
  title,
  description,
  children,
}: SettingItemProps) => {
  return (
    <>
      <Card
        className="border-medium border-foreground-300 w-full"
        classNames={{
          header: "p-4",
          body: "p-4",
        }}
        radius="sm"
        shadow="none"
      >
        <CardHeader className="bg-foreground-50">
          <div className="flex flex-col">
            <p className="text-medium font-semibold">{title}</p>
            <p className="text-default-500 text-xs">{description}</p>
          </div>
        </CardHeader>
        <Divider />
        <CardBody className="">{children}</CardBody>
      </Card>
    </>
  );
};
