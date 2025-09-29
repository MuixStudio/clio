export let apiUriPrefix = "";

if (
  process.env.NEXT_PUBLIC_API_PREFIX
) {
  apiUriPrefix = process.env.NEXT_PUBLIC_API_PREFIX;
} else {
  apiUriPrefix = "http://127.0.0.1:5020";
}

export const API_URI_PREFIX: string = apiUriPrefix;