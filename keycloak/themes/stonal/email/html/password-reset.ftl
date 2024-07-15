<#import "template.ftl" as layout>
<@layout.htmlEmailLayout ; section>
    <#if section = "text">
        ${msg("passwordResetBodyHtml", linkExpiration, realmName, user.firstName, user.lastName)?no_esc}
    </#if>
    <#if section = "linkText">
        ${msg("passwordResetLinkTextHtml")?no_esc}
    </#if>
    <#if section = "passwordSetText">
        ${msg("passwordSetText")}
    </#if>
    <#if section = "homePageLink">
        ${msg("homeLink")}
    </#if>
</@layout.htmlEmailLayout>
