export interface CoreKeymap {
  C00: string;
  C01: string;
  C02: string;
  C03: string;
  C04: string;
  C05: string;
  C06: string;
  C07: string;
  C08: string;
  C09: string;
  C10: string;
  C11: string;
  C12: string;
  C13: string;
  C14: string;
  C15: string;
  C16: string;
  C17: string;
  C18: string;
  C19: string;
  C20: string;
  C21: string;
  C22: string;
  C23: string;
  C24: string;
  C25: string;
  C26: string;
  C27: string;
  C28: string;
  C29: string;
}

export interface LeftKeymap {
  LB0: string;
  LB1: string;
  LB2: string;
}

export interface RightKeymap {
  RB0: string;
  RB1: string;
  RB2: string;
}

export interface ThumbKeymap {
  TB0: string;
  TB1: string;
  TB2: string;
  TB3: string;
  TB4: string;
  TB5: string;
}

export interface ExtraThumbKeymap {
  TE0: string;
  TE1: string;
}

export interface SpecialKeymap {
  SP0: string;
  SP1: string;
}

export interface NumberRowKeymap {
  N00: string;
  N01: string;
  N02: string;
  N03: string;
  N04: string;
  N05: string;
  N06: string;
  N07: string;
  N08: string;
  N09: string;
  N10: string;
  N11: string;
}

export type CompleteKeymap = CoreKeymap &
  LeftKeymap &
  RightKeymap &
  ThumbKeymap &
  ExtraThumbKeymap &
  SpecialKeymap &
  NumberRowKeymap;
