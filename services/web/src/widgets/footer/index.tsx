import React, { FC } from "react";
import FooterCard from "./card";
import FooterCardItem from "./cardItem";

type FooterProps = {
    absolute: boolean;
};

const Footer: FC<FooterProps> = ({ absolute }) => {
    let containerClass = "w-full";

    if (absolute) {
        containerClass += " absolute bottom-0 top-auto";
    }

    return (
        <div className={containerClass}>
            <style jsx>{`
                .footer-grid {
                    grid-template-columns: auto 1fr;
                }
            `}</style>
            <footer className="bg-blue-600 text-white font-sans md:px-16 md:py-8 py-4 px-8">
                <div className="md:grid footer-grid">
                    <div className="md:grid md:grid-rows-2 w-auto md:mr-16 mb-2">
                        <div>&copy; Jonas Schneider</div>
                        <a href="/" className="font-medium font-sans text-xl cursor-pointer">
                            Royalafg
                        </a>
                    </div>
                    <div className="md:grid footer-grid-content row-span-2 md:gap-2 md:grid-cols-3 md:justify-items-center">
                        <FooterCard title="legal">
                            <FooterCardItem href="/legal/terms">Terms & Conditions</FooterCardItem>
                            <FooterCardItem href="/legal/privacy">Privacy Statement</FooterCardItem>
                        </FooterCard>
                        <FooterCard title="Games">
                            <FooterCardItem href="/games">Game Selection</FooterCardItem>
                            <FooterCardItem href="/games/poker">Poker</FooterCardItem>
                            <FooterCardItem href="/games/pacman">Pacman</FooterCardItem>
                        </FooterCard>
                        <FooterCard title="legal">
                            <FooterCardItem href="/legal/terms">Terms & Conditions</FooterCardItem>
                            <FooterCardItem href="/legal/privacy">Privacy Statement</FooterCardItem>
                            <FooterCardItem>Hello</FooterCardItem>
                        </FooterCard>
                    </div>
                </div>
            </footer>
        </div>
    );
};

export default Footer;
