<template>
    <svg :width="width + 'em'" :height="height + 'em'" class="svg-icon" aria-hidden="true">
        <use :xlink:href="`#${type}`" />
    </svg>
</template>
<script lang="ts">
/**
 * 自定义 svg icon
 * @author LiQingSong
 * 使用说明：
 *   1、下载或制作svg文件，存放到/src/assets/iconsvg目录下
 *   2、项目根据/src/assets/iconsvg/svgo.yml配置自动压缩精简svg，也可以独立运行 npm run svgo 压缩精简svg
 *   3、使用Demo：
 *      import IconSvg from '@/components/IconSvg';
 *      import { defineComponent } from "vue";
 *      export default defineComponent({
 *          components: {
 *              IconSvg
 *          }
 *      })
 *      <IconSvg type="svg文件名" class="" style=""/>
 */
import { defineComponent } from "vue";
export default defineComponent({
    name: 'IconSvg',
    props: {
        type: {
            type: String,
            required: true
        },
      width: {
        type: Number,
        required: false,
        default: 1
      },
      height: {
        type: Number,
        required: false,
        default: 1
      }
    }
});
</script>
<style lang="less" scoped>
.svg-icon {
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}
</style>