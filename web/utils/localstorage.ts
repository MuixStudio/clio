export const removeAccessToken = () => {
  localStorage.removeItem("token");
  localStorage.removeItem("webapp_access_token");
};
