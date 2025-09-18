import React from "react";
import { Avatar } from "@heroui/avatar";
import { Spacer } from "@heroui/spacer";
import { Skeleton } from "@heroui/skeleton";
import { Tooltip } from "@heroui/tooltip";

export type MenuHeaderProps = {
  customHeaderContent?: undefined;
  isLoaded?: boolean;
  isCompact?: boolean;
  title?: string;
  description?: string;
};

const MenuHeader: React.FC<MenuHeaderProps> = ({ ...props }) => {
  if (props.customHeaderContent) {
    return <div className="sidebar-header">{props.customHeaderContent}</div>;
  } else {
    if (props.isCompact) {
      return (
        <Tooltip
          classNames={{ content: "max-w-36" }}
          closeDelay={100}
          content={props.description}
          placement="right-start"
        >
          <Skeleton className="rounded-lg" isLoaded={props.isLoaded}>
            <Avatar isBordered name={props.title} radius="md" />
          </Skeleton>
        </Tooltip>
      );
    } else {
      return (
        <div className="h-full p-3">
          <div className="flex h-full flex-col">
            <Skeleton className="rounded-lg" isLoaded={props.isLoaded}>
              <h1 className="text-medium text-default-900 min-h-[24px] font-bold">
                {props.title}
              </h1>
            </Skeleton>
            <Spacer y={1} />
            <Skeleton className="w-4/5 rounded-lg" isLoaded={props.isLoaded}>
              <h2 className="text-tiny text-default-500 min-h-[16px]">
                {props.description}
              </h2>
            </Skeleton>
          </div>
        </div>
      );
    }
  }
};

export default MenuHeader;
