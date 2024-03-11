const createResultItem = (name, no, imgSrc) => {
  const $div = document.createElement("div");
  $div.classList.add("result-item");
  const $a = document.createElement("a");
  $a.href = "https://zukan.pokemon.co.jp/detail/" + no;
  $a.target = "_blank";
  const $img = document.createElement("img");
  $img.src = imgSrc;
  const $span = document.createElement("span");
  $span.textContent = name;
  $a.appendChild($img);
  $a.appendChild($span);
  $div.appendChild($a);
  return $div;
};

(async () => {
  const res = await fetch("data/zukan.json");
  const data = await res.json();
  const $select = document.querySelector("#name-select");
  const $input = document.querySelector("#no-input");
  const noDist = {};
  data.data.forEach((p) => {
    const $option = document.createElement("option");
    if (Number(p.no) > 151) return;
    $option.value = p.no;
    $option.text = p.name;
    $select.appendChild($option);
    noDist[p.no] = p.name;
  });

  $input.addEventListener("change", (e) => {
    $select.value = e.target.value;
  });

  $input.addEventListener("input", (e) => {
    $select.value = e.target.value;
  });

  $input.addEventListener("blur", (e) => {
    const v = Number(e.target.value);
    if (v < 1) {
      e.target.value = 1;
    } else if (v > 151) {
      e.target.value = 151;
    }
    $select.value = e.target.value;
  });

  $select.addEventListener("change", (e) => {
    $input.value = e.target.value;
  });

  const $button = document.querySelector("button");
  const $result = document.querySelector("#result");
  $button.addEventListener("click", async () => {
    const path = goGetLongestChain(Number($select.value));
    if (path.length === 0) {
      $result.textContent =
        "このポケモンから始めると、しりとりは終了しません。";
      return;
    }

    $result.innerHTML = "";
    for (let i = 0; i < path.length; i++) {
      $result.appendChild(
        createResultItem(
          noDist[path[i]],
          path[i],
          `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/${path[i]}.png`
        )
      );
    }
  });
})();
