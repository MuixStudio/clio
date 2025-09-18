import React from "react";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { Avatar } from "@heroui/avatar";
import { Divider } from "@heroui/divider";

import {
  SettingCard,
  SettingCardContent,
  SettingCardGroup,
  SettingCardItem,
} from "../setting-card";

const Profile: React.FC = ({ className, ...props }: { className?: string }) => {
  return (
    <SettingCardGroup>
      {/* Avator */}
      <SettingCard description="your Avator." title="Avator">
        <SettingCardContent>
          <SettingCardItem title="Upload Avator">
            <div className="flex h-24 gap-6">
              <Avatar
                className="text-large h-24 w-24 shrink-0"
                src="https://i.pravatar.cc/150?u=a04258114e29026708c"
              />
              <div className="flex flex-col justify-center gap-2">
                <Button
                  className="bg-foreground-50 border-small max-w-16 shrink-0 font-medium"
                  disableAnimation={true}
                  size="sm"
                  variant="bordered"
                >
                  Upload
                </Button>
                <span className="text-default-500 text-xs">
                  zui jia chi cun{" "}
                  <span className="text-default-foreground font-medium">
                    192 * 192
                  </span>{" "}
                  px, max file size is 200 KiB. zui jia chi cun{" "}
                  <span className="text-default-foreground font-medium">
                    192 * 192
                  </span>{" "}
                  px, max file size is{" "}
                  <span className="text-default-foreground font-medium">
                    200
                  </span>{" "}
                  KiB.
                </span>
              </div>
            </div>
          </SettingCardItem>
        </SettingCardContent>
      </SettingCard>
      {/* Full name */}
      <SettingCard description="your profile information." title="Profile">
        <SettingCardContent>
          <SettingCardItem description="Your name" title="Name">
            <Input
              className="rounded-small mt-2 max-w-xl"
              classNames={{
                innerWrapper: "bg-foreground-50",
                inputWrapper: "bg-foreground-50",
              }}
              placeholder="username"
              radius="sm"
              variant="faded"
            />
          </SettingCardItem>
          <Divider />
          <SettingCardItem description="Your localtime" title="Localtime">
            <Input
              className="rounded-small mt-2 max-w-xl"
              classNames={{
                innerWrapper: "bg-foreground-50",
                inputWrapper: "bg-foreground-50",
              }}
              placeholder="e.g Kate Moore"
              radius="sm"
              variant="faded"
            />
          </SettingCardItem>
        </SettingCardContent>
      </SettingCard>
      <Button className="bg-default-foreground text-background" size="sm">
        Update Profile
      </Button>
    </SettingCardGroup>
  );
};

export default Profile;
