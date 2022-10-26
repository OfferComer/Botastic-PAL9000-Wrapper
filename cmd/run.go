package cmd

import (
	"context"
	"fmt"

	"github.com/pandodao/PAL9000/config"
	"github.com/pandodao/PAL9000/internal/discord"
	"github.com/pandodao/PAL9000/internal/mixin"
	"github.com/pandodao/PAL9000/internal/telegram"
	"github.com/pandodao/PAL9000/internal/wechat"
	"github.com/pandodao/PAL9000/service"
	"github.com/pandodao/PAL9000/store"
	"github.com/spf13/cobra"
	"golang.org/x/sync/e