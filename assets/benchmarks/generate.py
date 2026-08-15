#!/usr/bin/env python3
"""Generate benchmark SVGs from Go benchmark output.

Usage:
    python3 assets/benchmarks/generate.py \
        assets/benchmarks/raw/bench7.txt \
        assets/benchmarks/raw/bench-time.txt \
        assets/benchmarks/raw/bench-large-struct-parallel.txt
"""

from __future__ import annotations

import html
import re
import statistics
import sys
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path


HERE = Path(__file__).resolve().parent

GREEN = "#50c878"
BLUE = "#5aa9e6"
ORANGE = "#ef8354"
PINK = "#e76f9a"
PURPLE = "#b388eb"
YELLOW = "#f2c14e"


@dataclass(frozen=True)
class Result:
    ns: float
    mbps: float
    bytes: int
    allocs: int


def parse(path: Path, keep_cpu: bool = False) -> dict[tuple[str, int], Result]:
    samples: dict[tuple[str, int], list[Result]] = defaultdict(list)
    for line in path.read_text().splitlines():
        if not line.startswith("Benchmark"):
            continue
        fields = line.split()
        if len(fields) < 9 or fields[3] != "ns/op":
            continue
        name = fields[0]
        cpu = 1
        match = re.search(r"-(\d+)$", name)
        if match:
            cpu = int(match.group(1))
            name = name[: match.start()]
        key = (name, cpu if keep_cpu else 0)
        samples[key].append(
            Result(float(fields[2]), float(fields[4]), int(fields[6]), int(fields[8]))
        )

    medians: dict[tuple[str, int], Result] = {}
    for key, values in samples.items():
        medians[key] = Result(
            statistics.median(value.ns for value in values),
            statistics.median(value.mbps for value in values),
            round(statistics.median(value.bytes for value in values)),
            round(statistics.median(value.allocs for value in values)),
        )
    return medians


def latency(value: float) -> str:
    if value < 10:
        return f"{value:.3f}"
    if value < 100:
        return f"{value:.2f}"
    if value < 1000:
        return f"{value:.1f}"
    return f"{value:.0f}"


def allocs(value: int) -> str:
    return f"{value} alloc" if value == 1 else f"{value} allocs"


WORKLOADS = {
    "BenchmarkMarshalMapInt": "map[string]int",
    "BenchmarkMarshalMapIntSlice": "map[string][]int",
    "BenchmarkMarshalMapAny": "map[string]any",
    "BenchmarkMarshalIntSlice": "[]int",
    "BenchmarkMarshalFloat32": "float32",
    "BenchmarkMarshalFloat64": "float64",
    "BenchmarkMarshalTime": "time.Time",
    "BenchmarkMarshalTimeStruct": "struct with 10 time.Time fields",
    "BenchmarkMarshalStruct": "mixed struct",
    "BenchmarkMarshalStructSlice": "struct slice with metadata maps",
    "BenchmarkMarshalStructQuoted": "quoted struct fields",
    "BenchmarkMarshalOmits": "omitempty / omitzero",
    "BenchmarkMarshalOmitZero": "omitzero fields",
    "BenchmarkMarshalUTF8/valid_ascii": "UTF-8 validation: ASCII",
    "BenchmarkMarshalUTF8/valid_unicode": "UTF-8 validation: Unicode",
    "BenchmarkMarshalUTF8/invalid_middle": "UTF-8 validation: invalid byte",
    "BenchmarkTextMarshaler": "encoding.TextMarshaler",
    "BenchmarkJsonMarshaler": "json.Marshaler",
    "BenchmarkJsonIntMarshaler": "integer json.Marshaler",
}

VALUES = [
    "BenchmarkMarshalFloat32",
    "BenchmarkMarshalFloat64",
    "BenchmarkMarshalTime",
    "BenchmarkMarshalTimeStruct",
    "BenchmarkMarshalIntSlice",
    "BenchmarkMarshalStruct",
    "BenchmarkMarshalStructQuoted",
    "BenchmarkMarshalStructSlice",
]

MAPS_OPTIONS = [
    "BenchmarkMarshalMapInt",
    "BenchmarkMarshalMapIntSlice",
    "BenchmarkMarshalMapAny",
    "BenchmarkMarshalOmits",
    "BenchmarkMarshalOmitZero",
]

UTF8_INTERFACES = [
    "BenchmarkMarshalUTF8/valid_ascii",
    "BenchmarkMarshalUTF8/valid_unicode",
    "BenchmarkMarshalUTF8/invalid_middle",
    "BenchmarkTextMarshaler",
    "BenchmarkJsonMarshaler",
    "BenchmarkJsonIntMarshaler",
]

OVERVIEW = [
    "BenchmarkMarshalMapInt",
    "BenchmarkMarshalMapIntSlice",
    "BenchmarkMarshalMapAny",
    "BenchmarkMarshalIntSlice",
    "BenchmarkMarshalFloat32",
    "BenchmarkMarshalFloat64",
    "BenchmarkMarshalTime",
    "BenchmarkMarshalTimeStruct",
    "BenchmarkMarshalStruct",
    "BenchmarkMarshalStructSlice",
    "BenchmarkMarshalStructQuoted",
    "BenchmarkMarshalOmits",
    "BenchmarkMarshalOmitZero",
    "BenchmarkMarshalUTF8/valid_ascii",
    "BenchmarkMarshalUTF8/valid_unicode",
    "BenchmarkMarshalUTF8/invalid_middle",
    "BenchmarkTextMarshaler",
    "BenchmarkJsonMarshaler",
    "BenchmarkJsonIntMarshaler",
]


def result(data: dict[tuple[str, int], Result], base: str, suffix: str) -> Result:
    return data[(f"{base}/{suffix}", 0)]


def comparison_svg(
    data: dict[tuple[str, int], Result], title: str, workloads: list[str]
) -> str:
    height = 145 + (len(workloads) - 1) * 174 + 158 + 51
    out = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="1500" height="{height}" viewBox="0 0 1500 {height}" role="img" aria-labelledby="title desc">',
        f'  <title id="title">{html.escape(title)} — Benchmark 7</title>',
        '  <desc id="desc">Five-run median JSON marshal latency. Owned-output and reusable-output APIs are compared separately. Lower is better.</desc>',
        '  <style>text{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;fill:#d8dee9}.title{font-family:ui-sans-serif,system-ui,sans-serif;font-size:28px;font-weight:700}.subtitle{font-family:ui-sans-serif,system-ui,sans-serif;font-size:15px;fill:#aab4c3}.group{font-family:ui-sans-serif,system-ui,sans-serif;font-size:18px;font-weight:650}.contract{font-family:ui-sans-serif,system-ui,sans-serif;font-size:14px;font-weight:650;fill:#9eabbc}.series{font-size:13px;fill:#c4ccd8}.value{font-size:13px;font-weight:650}.panel{fill:#151c27}</style>',
        f'  <rect width="1500" height="{height}" rx="14" fill="#10151d"/>',
        f'  <text class="title" x="40" y="48">{html.escape(title)} — Benchmark 7</text>',
        '  <text class="subtitle" x="40" y="76">Go 1.26 · GOAMD64=v3 · JSON v2 + SIMD · Ryzen 9 9950X3D · five-run median · lower is better</text>',
        f'  <rect x="40" y="94" width="14" height="14" rx="3" fill="{GREEN}"/><text class="series" x="62" y="106">json-experiment</text>',
        f'  <rect x="225" y="94" width="14" height="14" rx="3" fill="{ORANGE}"/><text class="series" x="247" y="106">Go standard library</text>',
        f'  <rect x="445" y="94" width="14" height="14" rx="3" fill="{PURPLE}"/><text class="series" x="467" y="106">Sonic</text>',
        '  <text class="subtitle" x="1460" y="106" text-anchor="end">Each workload uses its own linear scale; labels show ns/op and allocations.</text>',
    ]

    for index, base in enumerate(workloads):
        y = 145 + index * 174
        owned = [
            ("json-experiment", result(data, base, "marshal"), GREEN),
            ("encoding/json", result(data, base, "encoding_json"), ORANGE),
            ("Sonic", result(data, base, "sonic_json"), PURPLE),
        ]
        reusable = [
            ("json-experiment", result(data, base, "marshal_append"), GREEN),
            ("json/v2", result(data, base, "encoding_json_v2_write"), ORANGE),
            ("Sonic", result(data, base, "sonic_encode_into"), PURPLE),
        ]
        out.extend(
            [
                f'  <g transform="translate(0,{y})">',
                '    <rect class="panel" x="25" y="0" width="1450" height="158" rx="10"/>',
                f'    <text class="group" x="45" y="27">{html.escape(WORKLOADS[base])}</text>',
                '    <text class="contract" x="45" y="53">Owned output</text>',
                '    <text class="contract" x="785" y="53">Reusable output</text>',
            ]
        )
        for contract, start_x, label_x in ((owned, 205, 193), (reusable, 945, 933)):
            maximum = max(item[1].ns for item in contract)
            for row, (name, value, color) in enumerate(contract):
                row_y = 68 + row * 25
                width = 360 * value.ns / maximum
                out.extend(
                    [
                        f'    <text class="series" x="{label_x}" y="{row_y + 13}" text-anchor="end">{html.escape(name)}</text>',
                        f'    <rect x="{start_x}" y="{row_y}" width="{width:.1f}" height="17" rx="4" fill="{color}"/>',
                        f'    <text class="value" x="{start_x + width + 8:.1f}" y="{row_y + 13}">{latency(value.ns)} ns · {allocs(value.allocs)}</text>',
                    ]
                )
        out.append("  </g>")
    out.append("</svg>")
    return "\n".join(out) + "\n"


def overview_svg(data: dict[tuple[str, int], Result]) -> str:
    height = 125 + (len(OVERVIEW) - 1) * 212 + 196 + 51
    series = [
        ("json-experiment MarshalAppend", "marshal_append", GREEN),
        ("json-experiment Marshal", "marshal", BLUE),
        ("encoding/json Marshal", "encoding_json", ORANGE),
        ("encoding/json/v2 MarshalWrite", "encoding_json_v2_write", PINK),
        ("sonic.ConfigFastest Marshal", "sonic_json", PURPLE),
        ("sonic EncodeInto", "sonic_encode_into", YELLOW),
    ]
    out = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="1500" height="{height}" viewBox="0 0 1500 {height}" role="img" aria-labelledby="title desc">',
        '  <title id="title">Complete JSON marshal comparison — Benchmark 7</title>',
        '  <desc id="desc">Five-run median JSON marshal latency across six APIs. Lower is better.</desc>',
        '  <style>text{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;fill:#d8dee9}.title{font-family:ui-sans-serif,system-ui,sans-serif;font-size:28px;font-weight:700}.subtitle{font-family:ui-sans-serif,system-ui,sans-serif;font-size:15px;fill:#aab4c3}.group{font-family:ui-sans-serif,system-ui,sans-serif;font-size:18px;font-weight:650}.series{font-size:13px;fill:#c4ccd8}.value{font-size:13px;font-weight:650}.panel{fill:#151c27}</style>',
        f'  <rect width="1500" height="{height}" rx="14" fill="#10151d"/>',
        '  <text class="title" x="40" y="48">Complete JSON marshal comparison — Benchmark 7</text>',
        '  <text class="subtitle" x="40" y="76">Go 1.26 · GOAMD64=v3 · JSON v2 + SIMD · Ryzen 9 9950X3D · five-run median · lower is better</text>',
        '  <text class="subtitle" x="1460" y="76" text-anchor="end">Each workload uses its own linear scale; labels show ns/op and allocations.</text>',
    ]
    for index, base in enumerate(OVERVIEW):
        y = 125 + index * 212
        values = [(name, result(data, base, suffix), color) for name, suffix, color in series]
        maximum = max(value.ns for _, value, _ in values)
        out.extend(
            [
                f'  <g transform="translate(0,{y})">',
                '    <rect class="panel" x="25" y="0" width="1450" height="196" rx="10"/>',
                f'    <text class="group" x="45" y="28">{html.escape(WORKLOADS[base])}</text>',
            ]
        )
        for row, (name, value, color) in enumerate(values):
            row_y = 43 + row * 25
            width = 720 * value.ns / maximum
            out.extend(
                [
                    f'    <text class="series" x="378" y="{row_y + 13}" text-anchor="end">{html.escape(name)}</text>',
                    f'    <rect x="390" y="{row_y}" width="{width:.1f}" height="17" rx="4" fill="{color}"/>',
                    f'    <text class="value" x="{398 + width:.1f}" y="{row_y + 13}">{latency(value.ns)} ns · {allocs(value.allocs)}</text>',
                ]
            )
        out.append("  </g>")
    out.append("</svg>")
    return "\n".join(out) + "\n"


def large_struct_svg(data: dict[tuple[str, int], Result]) -> str:
    methods = [
        ("MarshalAppend", "marshal_append", GREEN),
        ("Marshal", "marshal", GREEN),
        ("encoding/json", "encoding_json", ORANGE),
        ("json/v2 Write", "encoding_json_v2_write", ORANGE),
        ("Sonic Marshal", "sonic_json", PURPLE),
        ("Sonic EncodeInto", "sonic_encode_into", PURPLE),
    ]
    rows = sorted(
        ((label, result(data, "BenchmarkMarshalLargeStruct", suffix), color) for label, suffix, color in methods),
        key=lambda item: item[1].ns,
    )
    maximum = max(item[1].ns for item in rows)
    out = [
        '<svg xmlns="http://www.w3.org/2000/svg" width="1500" height="500" viewBox="0 0 1500 500" role="img" aria-labelledby="title desc">',
        '  <title id="title">Large nested struct JSON marshal benchmark</title>',
        '  <desc id="desc">Five-run median latency, throughput, and allocations for a 4,740-byte nested JSON payload on an AMD Ryzen 9 9950X3D. All APIs are ranked together by latency.</desc>',
        '  <style>text{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;fill:#d8dee9}.title{font-family:ui-sans-serif,system-ui,sans-serif;font-size:28px;font-weight:700}.subtitle{font-family:ui-sans-serif,system-ui,sans-serif;font-size:15px;fill:#aab4c3}.contract{font-family:ui-sans-serif,system-ui,sans-serif;font-size:18px;font-weight:650}.series{font-size:14px;fill:#c4ccd8}.value{font-size:14px;font-weight:650}.panel{fill:#151c27}.note{font-family:ui-sans-serif,system-ui,sans-serif;font-size:13px;fill:#8794a6}</style>',
        '  <rect width="1500" height="500" rx="14" fill="#10151d"/>',
        '  <text class="title" x="40" y="48">Large nested struct</text>',
        '  <text class="subtitle" x="40" y="76">Go 1.26 · GOAMD64=v3 · JSON v2 + SIMD · Ryzen 9 9950X3D · five-run median</text>',
        f'  <rect x="40" y="94" width="14" height="14" rx="3" fill="{GREEN}"/><text class="series" x="62" y="106">json-experiment</text>',
        f'  <rect x="225" y="94" width="14" height="14" rx="3" fill="{ORANGE}"/><text class="series" x="247" y="106">Go standard library</text>',
        f'  <rect x="445" y="94" width="14" height="14" rx="3" fill="{PURPLE}"/><text class="series" x="467" y="106">Sonic</text>',
        '  <text class="subtitle" x="1460" y="106" text-anchor="end">All APIs ranked by latency; lower is better.</text>',
        '  <g transform="translate(0,125)">',
        '    <rect class="panel" x="25" y="0" width="1450" height="278" rx="10"/>',
        '    <text class="contract" x="45" y="29">Fastest to slowest</text>',
    ]
    for row, (name, value, color) in enumerate(rows):
        row_y = 47 + row * 38
        width = 700 * value.ns / maximum
        out.extend(
            [
                f'    <text class="series" x="220" y="{row_y + 16}" text-anchor="end">{html.escape(name)}</text>',
                f'    <rect x="240" y="{row_y}" width="{width:.1f}" height="22" rx="5" fill="{color}"/>',
                f'    <text class="value" x="{250 + width:.1f}" y="{row_y + 16}">{value.ns / 1000:.3f} µs · {value.mbps:.0f} MB/s · {value.bytes} B · {allocs(value.allocs)}</text>',
            ]
        )
    out.extend(
        [
            "  </g>",
            '  <text class="note" x="40" y="477">4,740-byte encoded payload · nested structs, slices, maps, field options, ASCII, and Unicode</text>',
            "</svg>",
        ]
    )
    return "\n".join(out) + "\n"


def parallel_svg(data: dict[tuple[str, int], Result]) -> str:
    cpus = [1, 2, 4, 8, 16, 32]
    methods = [
        ("MarshalAppend", "marshal_append", "append", GREEN, 5),
        ("Marshal", "marshal", "marshal", GREEN, 4),
        ("Sonic EncodeInto", "sonic_encode_into", "sonic-into", PURPLE, 5),
        ("Sonic Marshal", "sonic_json", "sonic-marshal", PURPLE, 4),
    ]
    series = []
    for label, suffix, css, color, radius in methods:
        values = [data[(f"BenchmarkMarshalLargeStruct/parallel/{suffix}", cpu)] for cpu in cpus]
        series.append((label, css, color, radius, values))

    x_values = [170, 410, 650, 890, 1130, 1370]
    chart_bottom = 545
    scale = 13

    def y(value: float) -> float:
        return chart_bottom - value / 1000 * scale

    out = [
        '<svg xmlns="http://www.w3.org/2000/svg" width="1500" height="650" viewBox="0 0 1500 650" role="img" aria-labelledby="title desc">',
        '  <title id="title">Large nested struct parallel scaling</title>',
        '  <desc id="desc">Aggregate JSON encoding throughput from one to 32 logical CPUs for json-experiment and Sonic on an AMD Ryzen 9 9950X3D.</desc>',
        '  <style>text{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;fill:#d8dee9}.title{font-family:ui-sans-serif,system-ui,sans-serif;font-size:28px;font-weight:700}.subtitle{font-family:ui-sans-serif,system-ui,sans-serif;font-size:15px;fill:#aab4c3}.series{font-size:14px;fill:#c4ccd8}.tick{font-size:13px;fill:#8794a6}.axis{font-family:ui-sans-serif,system-ui,sans-serif;font-size:14px;font-weight:650;fill:#aab4c3}.value{font-size:13px;font-weight:700}.grid{stroke:#293344;stroke-width:1}.physical{stroke:#5d6a7c;stroke-width:1;stroke-dasharray:5 5}.append{fill:none;stroke:#50c878;stroke-width:5;stroke-linejoin:round;stroke-linecap:round}.marshal{fill:none;stroke:#50c878;stroke-width:4;stroke-dasharray:10 7;stroke-linejoin:round;stroke-linecap:round}.sonic-into{fill:none;stroke:#b388eb;stroke-width:5;stroke-linejoin:round;stroke-linecap:round}.sonic-marshal{fill:none;stroke:#b388eb;stroke-width:4;stroke-dasharray:10 7;stroke-linejoin:round;stroke-linecap:round}.point-green{fill:#50c878}.point-purple{fill:#b388eb}.panel{fill:#151c27}</style>',
        '  <rect width="1500" height="650" rx="14" fill="#10151d"/>',
        '  <text class="title" x="40" y="48">Large nested struct — parallel scaling</text>',
        '  <text class="subtitle" x="40" y="76">Go 1.26 · GOAMD64=v3 · JSON v2 + SIMD · Ryzen 9 9950X3D · three-run median</text>',
        '  <line x1="40" y1="102" x2="78" y2="102" class="append"/><text class="series" x="88" y="107">MarshalAppend</text>',
        '  <line x1="250" y1="102" x2="288" y2="102" class="marshal"/><text class="series" x="298" y="107">Marshal</text>',
        '  <line x1="405" y1="102" x2="443" y2="102" class="sonic-into"/><text class="series" x="453" y="107">Sonic EncodeInto</text>',
        '  <line x1="665" y1="102" x2="703" y2="102" class="sonic-marshal"/><text class="series" x="713" y="107">Sonic Marshal</text>',
        '  <text class="subtitle" x="1460" y="107" text-anchor="end">Aggregate throughput; higher is better.</text>',
        '  <rect class="panel" x="25" y="125" width="1450" height="465" rx="10"/>',
    ]
    for throughput in range(0, 31, 5):
        line_y = chart_bottom - throughput * scale
        out.append(f'  <line class="grid" x1="145" y1="{line_y}" x2="1405" y2="{line_y}"/><text class="tick" x="130" y="{line_y + 5}" text-anchor="end">{throughput}</text>')
    out.extend(
        [
            '  <text class="axis" transform="translate(58 350) rotate(-90)" text-anchor="middle">Aggregate throughput (GB/s)</text>',
            '  <line class="physical" x1="1130" y1="145" x2="1130" y2="555"/>',
            '  <text class="tick" x="1120" y="174" text-anchor="end">16 physical cores</text>',
        ]
    )
    for _, css, color, radius, values in series:
        points = " ".join(f"{x},{y(value.mbps):.1f}" for x, value in zip(x_values, values))
        out.append(f'  <polyline class="{css}" points="{points}"/>')
        circles = "".join(
            f'<circle cx="{x}" cy="{y(value.mbps):.1f}" r="{radius}"/>'
            for x, value in zip(x_values, values)
        )
        point_class = "point-green" if color == GREEN else "point-purple"
        out.append(f'  <g class="{point_class}">{circles}</g>')
    label_offsets = [-4, 4, -2, 5]
    for (_, _, color, _, values), offset in zip(series, label_offsets):
        value = values[-1]
        out.append(
            f'  <text class="value" x="1390" y="{y(value.mbps) + offset:.1f}" fill="{color}">{value.mbps / 1000:.2f}</text>'
        )
    for x, cpu in zip(x_values, cpus):
        out.append(f'  <text class="tick" x="{x}" y="574" text-anchor="middle">{cpu}</text>')
    out.extend(
        [
            '  <text class="axis" x="775" y="613" text-anchor="middle">GOMAXPROCS</text>',
            '  <text class="subtitle" x="40" y="635">4,740-byte payload · b.RunParallel · private reusable buffer per worker · solid = reusable output · dashed = owned output</text>',
            "</svg>",
        ]
    )
    return "\n".join(out) + "\n"


def main() -> None:
    if len(sys.argv) < 3:
        raise SystemExit(
            f"usage: {sys.argv[0]} BENCHMARK [BENCHMARK ...] PARALLEL_BENCHMARK"
        )
    benchmark = {}
    for path in sys.argv[1:-1]:
        benchmark.update(parse(Path(path)))
    parallel = parse(Path(sys.argv[-1]), keep_cpu=True)

    outputs = {
        "benchmark7-values.svg": comparison_svg(benchmark, "Values and structs", VALUES),
        "benchmark7-maps-options.svg": comparison_svg(benchmark, "Maps and field options", MAPS_OPTIONS),
        "benchmark7-utf8-interfaces.svg": comparison_svg(benchmark, "UTF-8 and marshaling interfaces", UTF8_INTERFACES),
        "benchmark7.svg": overview_svg(benchmark),
        "large-struct7.svg": large_struct_svg(benchmark),
        "large-struct-parallel7.svg": parallel_svg(parallel),
    }
    for name, contents in outputs.items():
        (HERE / name).write_text(contents)


if __name__ == "__main__":
    main()
