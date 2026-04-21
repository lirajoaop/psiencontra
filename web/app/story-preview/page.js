import { Fraunces } from "next/font/google";
import StoryPreviewClient from "./StoryPreviewClient";

const fraunces = Fraunces({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  style: ["normal", "italic"],
  variable: "--font-fraunces",
  display: "swap",
});

export const metadata = {
  title: "Preview · Story card",
  robots: { index: false, follow: false },
};

export default function Page() {
  return (
    <div className={fraunces.variable}>
      <StoryPreviewClient />
    </div>
  );
}
