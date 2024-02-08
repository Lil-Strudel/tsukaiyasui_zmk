import {
  qwertyCoreKeymap,
  leftKeymap,
  rightKeymap,
  thumbKeymap,
  extraThumbKeymap,
  specialKeymap,
  numberRowKeymap,
} from "./keymaps";

import { CompleteKeymap } from "./types";
import { readFileSync } from "node:fs";

interface MainProps {
  shield: "corne" | "lily58";
  layout: "qwerty";
}
function main({ shield, layout }: MainProps) {
  const adapter = readFileSync(`./adapters/${shield}.tsukaiyasui`, "utf-8");

  const coreKeymap = (() => {
    switch (layout) {
      case "qwerty":
        return qwertyCoreKeymap;
    }
  })();

  const completeKeymap: CompleteKeymap = {
    ...coreKeymap,
    ...leftKeymap,
    ...rightKeymap,
    ...thumbKeymap,
    ...extraThumbKeymap,
    ...specialKeymap,
    ...numberRowKeymap,
  };

  const rows = adapter.split("\n").filter((line) => line !== "");
  const splitParts = rows.map((row) => row.split(" "));

  const transformedParts = splitParts.map((row) =>
    row.map((key) =>
      key === "___" ? "" : completeKeymap[key as keyof CompleteKeymap],
    ),
  );
}

main({
  shield: "corne",
  layout: "qwerty",
});
