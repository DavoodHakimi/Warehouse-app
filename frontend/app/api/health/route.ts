import { NextResponse } from "next/server";

// Liveness probe for the Docker HEALTHCHECK and orchestrators.
export function GET() {
  return NextResponse.json({ status: "ok" }, { status: 200 });
}
