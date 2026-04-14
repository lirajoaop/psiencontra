export const GLOSSARY = {
  SUAS: "Sistema Único de Assistência Social — rede pública que organiza a proteção social no Brasil (inclui CRAS e CREAS).",
  CRAS: "Centro de Referência de Assistência Social — unidade do SUAS voltada à proteção social básica.",
  CREAS: "Centro de Referência Especializado de Assistência Social — unidade do SUAS voltada a famílias e indivíduos em situação de risco ou violação de direitos.",
  UBS: "Unidade Básica de Saúde — porta de entrada do SUS na atenção primária.",
  SUS: "Sistema Único de Saúde — sistema público e universal de saúde do Brasil.",
  ONG: "Organização Não Governamental.",
  AVC: "Acidente Vascular Cerebral.",
  TCE: "Traumatismo Cranioencefálico.",
  TDAH: "Transtorno do Déficit de Atenção e Hiperatividade.",
};

export const GLOSSARY_REGEX = new RegExp(
  `\\b(${Object.keys(GLOSSARY).join("|")})\\b`,
  "g"
);
