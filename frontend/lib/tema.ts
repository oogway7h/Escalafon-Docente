export const CLAVE_TEMA = "escalafon:tema";

// Se ejecuta en <head> antes de pintar: aplica el tema guardado o, si no hay, el del sistema.
export const SCRIPT_TEMA = `try{var t=localStorage.getItem("${CLAVE_TEMA}");if(t==="oscuro"||(!t&&matchMedia("(prefers-color-scheme: dark)").matches))document.documentElement.classList.add("dark")}catch(e){}`;
