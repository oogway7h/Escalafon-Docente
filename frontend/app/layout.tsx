import type { Metadata } from "next";
import { Manrope, Source_Sans_3 } from "next/font/google";
import Script from "next/script";
import "./globals.css";
import { SCRIPT_TEMA } from "@/lib/tema";

const manrope = Manrope({
  variable: "--font-manrope",
  subsets: ["latin"],
  weight: ["500", "700", "800"],
});

const sourceSans = Source_Sans_3({
  variable: "--font-source-sans",
  subsets: ["latin"],
  weight: ["400", "600", "700"],
});

export const metadata: Metadata = {
  title: "Escalafón Docente · UAGRM",
  description:
    "Sistema de Gestión de Flujos de Trabajo y Rastreo Documental del Departamento de Escalafón Docente, UAGRM.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    // suppressHydrationWarning: SCRIPT_TEMA agrega la clase "dark" antes de hidratar.
    <html lang="es" className={`${manrope.variable} ${sourceSans.variable} h-full`} suppressHydrationWarning>
      <body className="h-full font-[family-name:var(--font-body)] antialiased">
        {/* beforeInteractive: Next lo inyecta en el HTML inicial y corre antes de hidratar (sin parpadeo de tema).
            Un <script> crudo en el layout hace que React avise que los scripts dentro de componentes no se ejecutan en el cliente. */}
        <Script id="tema" strategy="beforeInteractive" dangerouslySetInnerHTML={{ __html: SCRIPT_TEMA }} />
        {children}
      </body>
    </html>
  );
}
