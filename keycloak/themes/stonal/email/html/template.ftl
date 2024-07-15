<#macro htmlEmailLayout>
    <html>
    <head>
        <style type="text/css" media="all">

            @font-face {
                font-family: 'Open Sans';
                font-style: normal;
                font-weight: 400;
                font-stretch: 100%;
                src: url(https://fonts.gstatic.com/s/opensans/v34/memSYaGs126MiZpBA-UvWbX2vVnXBbObj2OVZyOOSr4dVJWUgsjZ0B4gaVIUwaEQbjA.woff2) format('woff2');
                unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
            }

            * {
                font-family: 'Open Sans', 'Helvetica Neue', sans-serif !important;
                color: #333333;
                background-color: #F8F8F8;
            }

            p, span, small, a, strong {
                background-color: white !important;
            }

            .custom-background {
                background-color: #F8F8F8 !important;
            }

            .global {
                font-family: 'Open Sans', 'Helvetica Neue', sans-serif !important;
                padding: 15px 0;
                margin: auto;
                width: 75%;
                background-color: #F8F8F8;
                color: #333333;
            }

            a {
                color: #6468DE !important;
                text-decoration: none;
            }

            strong {
                color: #1EDC78;
            }

            img {
                width: 25%;
                height: auto;
            }

            .card {
                border: 1px solid lightgray;
                padding: 15px;
                background-color: white !important;
            }

            .password_reset_link {
                background-color: #142D46 !important;
                color: white !important;
                border-radius: 5px;
                border: none;
                padding: 10px;
                margin: 10px 0;
                text-align: center;
                text-decoration: none;
                display: inline-block;
                font-size: 15px;
                font-weight: bold;
            }

            .center {
                text-align: center;
                margin: 15px 0;
                background-color: white !important;
            }
        </style>
    </head>
    <body class="global">
    <div>
        <div style="padding-bottom: 15px;">
            <img src="https://stonal.com/wp-content/uploads/2022/03/stonal-logo-horizontal-1.png"
                 alt="stonal_horizontal_logo">
        </div>
        <div class="card">
            <#nested "text">

            <div class="center">
                <a class="password_reset_link" href="${link}"><#nested "linkText"></a>
            </div>
            <p style="background-color: white">
                <span><#nested "passwordSetText"> </span>
                <#assign homePage><#nested "homePageLink"></#assign>
                <a class="home_link" href="${homePage}" target="_blank">this link</a>.
            </p>
            <p style="background-color: white">
                <span>Regards,</span><br/>
                <span>Team</span>
                <a href="http://www.stonal.com" target="_blank">
                    <strong style=" color: #1EDC78;">Stonal</strong>
                </a>
            </p>
        </div>
        <p class="custom-background">
            <small class="custom-background"><strong class="custom-background">Stonal</strong> | 28 cours Albert 1er,
                75008 Paris | T. 01 76 54 84 37</small>
        </p>
    </div>
    </body>
    </html>
</#macro>
