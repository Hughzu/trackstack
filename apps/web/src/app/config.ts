export type AppDomain = {
  id: string;
  label: string;
  href: string;
};

export type AppProfile = {
  initials?: string;
  imageUrl?: string;
  href?: string;
  ariaLabel?: string;
};

export type AppConfig = {
  appName: string;
  domains: AppDomain[];
  profile?: AppProfile;
};

export const appConfig: AppConfig = {
  appName: "TrackStack",
  domains: [
    { id: "home", label: "Overview", href: "/" },
    { id: "expenses", label: "Expenses", href: "/expenses" },
    { id: "calories", label: "Health", href: "/calories" },
    { id: "heat", label: "Heat", href: "/heat" }
  ],
  profile: {
    initials: "HZ",
    ariaLabel: "Open profile"
  }
};
