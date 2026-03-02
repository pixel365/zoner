#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./epp-auth-test.sh [ok|invalid|empty] [host] [port]

Examples:
  ./epp-auth-test.sh ok
  ./epp-auth-test.sh invalid localhost 7000
  ./epp-auth-test.sh empty

Environment overrides:
  EPP_USER        default: demo
  EPP_PASS        default: demo-pass
  EPP_BAD_PASS    default: wrong-pass
  EPP_OBJ_URI     default: urn:ietf:params:xml:ns:domain-1.0
  EPP_TIMEOUT_SEC default: 5
EOF
}

scenario="${1:-ok}"
host="${2:-localhost}"
port="${3:-7000}"

if [[ "${scenario}" != "ok" && "${scenario}" != "invalid" && "${scenario}" != "empty" ]]; then
  usage
  exit 1
fi

user="${EPP_USER:-demo}"
pass_ok="${EPP_PASS:-demo-pass}"
pass_bad="${EPP_BAD_PASS:-wrong-pass}"
obj_uri="${EPP_OBJ_URI:-urn:ietf:params:xml:ns:domain-1.0}"
timeout_sec="${EPP_TIMEOUT_SEC:-5}"

case "${scenario}" in
  ok)
    pass="${pass_ok}"
    tr_id="manual-login-ok"
    ;;
  invalid)
    pass="${pass_bad}"
    tr_id="manual-login-invalid"
    ;;
  empty)
    pass=""
    tr_id="manual-login-empty"
    ;;
esac

build_login_xml() {
  local clid="$1"
  local pw="$2"
  local obj="$3"
  local trid="$4"

  cat <<EOF
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>${clid}</clID><pw>${pw}</pw><options><version>1.0</version><lang>en</lang></options><svcs><objURI>${obj}</objURI></svcs></login><clTRID>${trid}</clTRID></command></epp>
EOF
}

epp_send_frame() {
  local payload="$1"
  local len
  len=$((4 + ${#payload}))
  { printf '%08x' "${len}" | xxd -r -p; printf '%s' "${payload}"; }
}

epp_read_frames() {
  perl -e '
    binmode(STDIN);
    while (read(STDIN, my $h, 4) == 4) {
      my $n = unpack("N", $h) - 4;
      last if $n < 0;
      read(STDIN, my $p, $n) == $n or last;
      print $p, "\n\n";
    }
  '
}

xml_payload="$(build_login_xml "${user}" "${pass}" "${obj_uri}" "${tr_id}")"

echo "Scenario: ${scenario}"
echo "Target: ${host}:${port}"
echo "User: ${user}"
if [[ "${scenario}" == "empty" ]]; then
  echo "Password: <empty>"
elif [[ "${scenario}" == "invalid" ]]; then
  echo "Password: <invalid>"
else
  echo "Password: <provided>"
fi
echo
echo "Expected result code:"
if [[ "${scenario}" == "ok" ]]; then
  echo "  1000 (Command completed successfully)"
elif [[ "${scenario}" == "invalid" || "${scenario}" == "empty" ]]; then
  echo "  2200 (Authentication error)"
fi
echo "  2400 indicates internal server error"
echo
echo "Response frames:"
echo "----------------------------------------"

{ epp_send_frame "${xml_payload}"; sleep 0.2; } \
  | timeout "${timeout_sec}s" openssl s_client -connect "${host}:${port}" -quiet \
  | epp_read_frames

echo "----------------------------------------"
