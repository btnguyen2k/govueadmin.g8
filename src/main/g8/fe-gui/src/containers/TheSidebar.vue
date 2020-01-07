<template>
    <CSidebar
            fixed
            :minimize="minimize"
            :show.sync="show"
    >
        <div class="c-sidebar-brand">
            <a href="/">
                <span class="c-sidebar-brand-full" style="color: #fff; font-weight: bolder; font-size: x-large">GoVueAdmin</span>
                <span class="c-sidebar-brand-minimized" style="color: #fff; font-weight: bolder; font-size: large">GVA</span>
            </a>
        </div>
        <!--
        <div class="c-sidebar-brand">
            <a rel="noopener" href="https://coreui.io/" target="_blank" class="">
                <img
                        width="118" height="46" alt="Logo" src="img/brand/coreui-base-white.svg"
                        class="c-sidebar-brand-full">
                <img width="118" height="46" alt="Logo"
                     src="img/brand/coreui-signet-white.svg"
                     class="c-sidebar-brand-minimized">
            </a>
        </div>
        -->
        <!--
        <CSidebarBrand
                :imgFull="{ width: 118, height: 46, alt: 'Logo', src: 'img/brand/coreui-base-white.svg'}"
                :imgMinimized="{ width: 118, height: 46, alt: 'Logo', src: 'img/brand/coreui-signet-white.svg'}"
                :wrappedInLink="{ href: 'https://coreui.io/', target: '_blank'}"
        />
        -->
        <CRenderFunction flat :content-to-render="nav"/>
        <CSidebarMinimizer
                class="d-md-down-none"
                @click.native="minimize = !minimize"
        />
    </CSidebar>
</template>

<script>
    import nav from './_nav'

    export default {
        name: 'TheSidebar',
        data() {
            return {
                minimize: false,
                nav,
                show: 'responsive'
            }
        },
        mounted() {
            this.$root.$on('toggle-sidebar', () => {
                const sidebarOpened = this.show === true || this.show === 'responsive'
                this.show = sidebarOpened ? false : 'responsive'
            })
            this.$root.$on('toggle-sidebar-mobile', () => {
                const sidebarClosed = this.show === 'responsive' || this.show === false
                this.show = sidebarClosed ? true : 'responsive'
            })
        }
    }
</script>
