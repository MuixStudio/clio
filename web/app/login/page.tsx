"use client";

import React from "react";
import { Button, Divider, Input, Link } from "@heroui/react";
import { Icon } from "@iconify/react";

import { Boxes } from "@/components/bg/bg";
import { Test } from "@/service/test";

export default function Component() {
  const [isVisible, setIsVisible] = React.useState(false);
  const [isConfirmVisible, setIsConfirmVisible] = React.useState(false);

  const toggleVisibility = () => setIsVisible(!isVisible);
  const toggleConfirmVisibility = () => setIsConfirmVisible(!isConfirmVisible);

  const Login = async () => {
    const res = await Test({ username: "quanli.wang", password: "happy" })
      console.log(res.name);
  };

  return (
    <div className="h-full relative w-full overflow-hidden bg-slate-900 flex flex-col items-center justify-center">
      <Boxes />
      <div className="z-20 rounded-large bg-content1 shadow-small flex w-full max-w-sm flex-col gap-4 px-8 pt-6 pb-10">
        <p className="pb-2 text-xl font-medium">登录</p>
        <form
          className="flex flex-col gap-3"
          onSubmit={(e) => e.preventDefault()}
        >
          <Input
            isRequired
            label="用户名"
            name="username"
            placeholder="输入你的用户名"
            type="text"
            variant="bordered"
          />
          <Input
            isRequired
            endContent={
              <button type="button" onClick={toggleVisibility}>
                {isVisible ? (
                  <Icon
                    className="text-default-400 pointer-events-none text-2xl"
                    icon="solar:eye-closed-linear"
                  />
                ) : (
                  <Icon
                    className="text-default-400 pointer-events-none text-2xl"
                    icon="solar:eye-bold"
                  />
                )}
              </button>
            }
            label="密码"
            name="password"
            placeholder="输入你的密码"
            type={isVisible ? "text" : "password"}
            variant="bordered"
          />
          <Button color="primary" type="submit" onClick={Login}>
            登录
          </Button>
        </form>
        <div className="flex items-center gap-4 py-2">
          <Divider className="flex-1" />
          <p className="text-tiny text-default-500 shrink-0">OR</p>
          <Divider className="flex-1" />
        </div>
        <div className="flex flex-col gap-2">
          <Button
            startContent={<Icon icon="flat-color-icons:google" width={24} />}
            variant="bordered"
          >
            使用Google登录
          </Button>
          <Button
            startContent={
              <Icon className="text-default-500" icon="fe:github" width={24} />
            }
            variant="bordered"
          >
            使用Github登录
          </Button>
        </div>
        <p className="text-small text-center">
          没有账户?&nbsp;点击
          <Link className="text-small font-medium" href="#">
            注册
          </Link>
        </p>

        <p className="text-xs text-center">
          使用即代表您同意我们的&nbsp;
          <Link className="text-xs font-black" href="#">
            使用协议
          </Link>
          &&nbsp;
          <Link className="text-xs font-black" href="#">
            隐私政策
          </Link>
        </p>
      </div>
    </div>
  );
}
