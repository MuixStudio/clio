import { Button } from "@heroui/react";
import { Icon } from "@iconify/react";
import React from "react";


export default function SigninWithSSO() {
  const click = () => {
  };

  return (
    <Button
      disableAnimation={true}
      startContent={<Icon icon="flat-color-icons:google" width={24} />}
      variant="bordered"
      onPress={click}
    >
      <span className="font-semibold">Continue with Google</span>
    </Button>
  );
}
