import Insignia from "@/components/common/Insignia";
import { ETIQUETA_ESTADO, TONO_ESTADO, type EstadoTramite } from "@/lib/tramites";

export default function EstadoInsignia({ estado }: { estado: EstadoTramite }) {
  return <Insignia tono={TONO_ESTADO[estado]}>{ETIQUETA_ESTADO[estado]}</Insignia>;
}
