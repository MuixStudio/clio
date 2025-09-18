import React from 'react'
const Menu2: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <>
    <header className="rounded-medium border-small border-divider flex items-center gap-3 p-4">
      <h2 className="text-medium text-default-700 font-medium">Overview</h2>
    </header>
    <main className="mt-4 h-full w-full overflow-visible">{children}</main>
  </>
)
export default Menu2
