import React from "react";

export type MenuFooterProps = { customFooterContent?: React.ReactNode };

const MenuFooter: React.FC<MenuFooterProps> = ({ customFooterContent }) => {
  if (customFooterContent)
    return <div className="sidebar-footer">{customFooterContent}</div>;
};

export default MenuFooter;
