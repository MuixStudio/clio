"use client";

import { usePathname } from "next/navigation";
import { useMemo } from "react";

/**
 * 自定义钩子：判断菜单项是否激活
 *
 * 根据当前路由路径判断菜单项是否应该被高亮显示
 *
 * @param itemUrl - 菜单项的 URL
 * @param exact - 是否精确匹配（默认 false，使用前缀匹配）
 * @returns {boolean} 菜单项是否激活
 */
export function useIsActive(itemUrl: string, exact: boolean = false): boolean {
  const pathname = usePathname();

  return useMemo(() => {
    if (exact) {
      // 精确匹配：路径完全相同
      return pathname === itemUrl;
    } else {
      // 前缀匹配：当前路径以菜单项 URL 开头
      // 处理根路径的特殊情况
      if (itemUrl === "/") {
        return pathname === "/";
      }

      // require a "/" boundary so "/cmdb" doesn't match "/cmdb/models"
      return pathname === itemUrl || pathname.startsWith(itemUrl + "/");
    }
  }, [pathname, itemUrl, exact]);
}

/**
 * 自定义钩子：获取当前激活的路径
 *
 * @returns {string} 当前路由路径
 */
export function useActivePath(): string {
  return usePathname();
}
