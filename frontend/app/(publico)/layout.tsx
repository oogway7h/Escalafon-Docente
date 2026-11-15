import PublicoShell from "@/components/layout/PublicoShell";

export default function PublicoLayout({ children }: { children: React.ReactNode }) {
  return <PublicoShell>{children}</PublicoShell>;
}
