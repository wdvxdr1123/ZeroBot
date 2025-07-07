package napcat

import (
	"fmt"
	"strconv"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	zero.OnCommand("napcat_test").Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("== NapCat API Test Start =="))

		uid := strconv.Itoa(int(ctx.Event.UserID))
		gid := strconv.Itoa(int(ctx.Event.GroupID))
		// 1. ArkSharePeer 如果uid或gid不存在，能把napcat干下线
		arkJSON := ctx.ArkSharePeer(uid, gid)
		ctx.SendChain(message.Text("ArkSharePeer Result:\n", arkJSON))
		time.Sleep(500 * time.Millisecond)

		// 2. ArkShareGroup
		groupArk := ctx.ArkShareGroup(gid)
		ctx.SendChain(message.Text("ArkShareGroup Result:\n", groupArk))
		time.Sleep(500 * time.Millisecond)

		// 3. GetRobotUinRange
		start, end := ctx.GetRobotUinRange()
		ctx.SendChain(message.Text(fmt.Sprintf("GetRobotUinRange Result: start=%d, end=%d", start, end)))
		time.Sleep(500 * time.Millisecond)

		// 4. TranslateEn2Zh 卡死了
		// translatedWords := ctx.TranslateEn2Zh([]string{"hello", "world"})
		// ctx.SendChain(message.Text(fmt.Sprintf("TranslateEn2Zh Result: %v", translatedWords)))
		// time.Sleep(500 * time.Millisecond)

		// 5. SendForwardMessage
		messageID, resID := ctx.SendForwardMessage("group", 0, ctx.Event.GroupID, message.Message{
			message.CustomNode("椛椛", ctx.Event.SelfID, "这是一条自定义信息"),
		})
		ctx.SendChain(message.Text(fmt.Sprintf("SendForwardMessage Result: messageID=%d, resID=%s", messageID, resID)))
		time.Sleep(500 * time.Millisecond)

		// 6. GetAIRecord
		aiRecord := ctx.GetAIRecord("lucy-voice-suxinjiejie", ctx.Event.GroupID, "这是一段测试语音")
		ctx.SendChain(message.Text("GetAIRecord Result (base64):\n", aiRecord))
		time.Sleep(500 * time.Millisecond)

		// 7. SendGroupAIRecord
		aiMessageID := ctx.SendGroupAIRecord("lucy-voice-suxinjiejie", ctx.Event.GroupID, "这是一段测试语音")
		ctx.SendChain(message.Text("SendGroupAIRecord Result (messageID):\n", aiMessageID))
		time.Sleep(500 * time.Millisecond)

		ctx.SendChain(message.Text("== NapCat API Test End =="))
	})
}
