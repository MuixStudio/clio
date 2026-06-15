export const removeSession = () => {
  localStorage.removeItem("session");
  localStorage.removeItem("webapp_access_token");
};
