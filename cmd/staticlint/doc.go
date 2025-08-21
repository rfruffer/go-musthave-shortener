/*
Package main — multichecker для github.com/rfruffer/go-musthave-shortener.

Запуск из корня проекта:

	go run ./cmd/staticlint ./...

или:

	go build -o staticlint ./cmd/staticlint && ./staticlint ./...

Состав анализаторов:

 1. Стандартные анализаторы (golang.org/x/tools/go/analysis/passes):
    assign, atomic, bools, buildssa, composite, copylock, deepequalerrors,
    errorsas, fieldalignment, httpresponse, ifaceassert, loopclosure,
    nilfunc, nilness, printf, shadow, stdmethods, structtag, tests,
    unreachable, unsafeptr, unusedresult.
    Назначение: поиск типичных проблем — неправильные форматы printf,
    возможные nil-ошибки, тени переменных, неиспользуемые результаты и пр.

2) Staticcheck (honnef.co/go/tools):
  - Все SA* (suspicious) — потенциальные баги и подозрительный код.
  - S* (simple): S1000/S1002 — упрощения и микро-оптимизации кода.
  - ST* (stylecheck): ST1000 — требования к комментариям пакетов.
    Назначение: повышение надёжности, читаемости и единообразия кода.

3) Публичные анализаторы:

  - bodyclose — проверяет корректное закрытие http.Response.Body.

  - nilerr — находит случаи несоответствия err и возвращаемых значений
    (например, return nil при err != nil).

    4. Кастомный анализатор (noosexitinmain):
    Запрещает прямой вызов os.Exit в функции main пакета main.
    Мотивация: выносить завершение в обработку ошибок (например, run() error),
    а в main() — логировать/завершать без прямого os.Exit (log.Fatal допустим).

Требования задания:
- Исходный код проекта должен проходить проверку без ошибок.
- При необходимости переписать main так, чтобы не было прямого os.Exit.

Пример типового каркаса main:

	func main() {
	    if err := run(); err != nil {
	        log.Fatalf("fatal: %v", err) // допустимо: не прямой os.Exit
	    }
	}
*/
package main
