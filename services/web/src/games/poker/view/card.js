import { rW, rH, isMobile } from "./utils";
import * as PIXI from "pixi.js";

export const CARDWIDTH = 80;
export const CARDHEIGHT = 130;
const emptyCardText = "";

class Card extends PIXI.Container {
    constructor(id, card) {
        super();

        this.card = {
            value: -1,
            color: -1
        };

        if (isMobile()) {
            this.sprite = new PIXI.Text(emptyCardText, {
                fontSize: 10
            });
            this.sprite.resolution = 2;
            this.addChild(this.sprite);
        } else {
            this.sprite = new PIXI.Sprite(PIXI.Texture.from("back.png"));
            this.addChild(this.sprite);
            let rect = new PIXI.Graphics();
            rect.lineStyle(1, 0x000000, 1);
            rect.drawRoundedRect(0, 0, rW(CARDWIDTH), rH(CARDHEIGHT), 5);
            this.addChild(rect);
        }
        this.update(card);
    }

    update(card) {
        this.card = {
            ...this.card,
            ...card
        };
        if (isMobile()) {
            let uni;
            switch (this.card.color) {
                case 0:
                    uni = "	♣";
                    break;
                case 1:
                    uni = "♦";
                    break;
                case 2:
                    uni = "♥";
                    break;
                case 3:
                    uni = "♠";
                    break;
                default:
                    uni = "";
                    break;
            }
            if (this.card.value === -1 || this.card.color === -1) {
                this.sprite.text = emptyCardText;
            } else {
                this.sprite.text = this.card.value + uni;
            }
        } else {
            if (this.card.color < 0 || this.card.value < 0) {
                this.sprite.texture = PIXI.Texture.from("back.png");
            } else {
                this.sprite.texture = PIXI.Texture.from(
                    `${this.card.value}_${this.card.color}.png`
                );
            }
            this.sprite.texture.baseTexture.scaleMode = PIXI.SCALE_MODES.NEAREST;
            this.sprite.width = rW(CARDWIDTH);
            this.sprite.height = rH(CARDHEIGHT);
        }
    }
}

export { Card };