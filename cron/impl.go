package cron

// parse ...
func parse(spec string) (crontabIx, error) {
  return parser.Parse(spec)
}
