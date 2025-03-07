// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+y9+XLcOJIw/ioY/n5/2N11SbZ3pxWxMStb6m7P+tBK8nbEWgoJIlEqjEmCTYCSy249",
	"1vcC35N9gZMgCbJYpVLpMCYm3CpcTCQSeQCZie9BSJKMpChlNNj5HtBwhhIo/nwNwy9F9iZHkKH3JEIx",
	"L8xykqGcYSSaRJDBC0iR+BvRMMcZwyQNdoLjGQK6FrAZZIDNELgQQ4ILFJP0kgJGgkGQwK/vUHrJZsHO",
	"v70YBAlO9c+tQZBBxlDOB/wMh99O+T+T4S+nPwWDgM0zFOwElOU4vQwGwdchgRkehiRClygdoq8sh0MG",
	"LwWgIZsHOyW4g2AWxtUCNqXRF7voRgyZRygPdrb532w6xBFKGZ5iXsbyAqniFCao1jWGFygWn4ZRhDlS",
	"YHxQwV05tWfPPu8O/1dN7fPQ/H02Ov3p+T+suueNad8Manj/RFE+jNAUpygCEggAGYPhDEWAEbEKOaKk",
	"yEO1LiFMwQUCBUURmJIcTHHMkMLp7WAkF/9CIeu3NApfemHMT70sqsBelBdV7JctSH4JU/wNSpS4KNNu",
	"8WCoswK2RkStUKOjUmwjZdKLUuvds5yIpXIiS1U+GDxpYDWKyt8aO7rERsxWL8SYnjeDIEd/FjhHUbDz",
	"2WYUtRXRPU7rZH8zUBzU884leSdOMpKzI8GjOAD/f46mwU7w/41LUTVWcmr81morEe2Z7yaY78s25isL",
	"XCTNawCZWsRcneRnPoOf/6oS6/PVqFVAoeejfujZiJ8uQdJNobqbly5eutxWunDhItnA/6CcthLSlazU",
	"m0b3GYE/ZigFNEMh/1QEcAogOD/4dHwOuNBClIEMzmMCowHAaYRDyBAtMSzHQYDOSBFHgv9kEWQoGgCY",
	"RgBTyY8u5qI1nVOGEs6+LguYRwBeQpxSBkKShkWeo5Sp7nS0LMY1RocVXP7bzSCgDLKCLmL+Ur4eibaK",
	"+bfLbcUFlhbfByTG4fw9phRFsmRZWT4t4ng+/LOAsVwuI9sFR+Rs/3qGwxmAmvavIQWJ+KBcBExBJqDY",
	"sOjeapfRCaIUXjomvAtmRQJTkCMYwYsYAdVSEyJOL0GEGMRcNl6QQpKknq3a+2uYtAZQz7n8raesSzps",
	"CquJhPAYt0k3hhMl03stI/oKkyzm39iebL8cTraGk63jyWRH/P9/g0EwJXkCmcQ6GvLRnXi4JMNVkCOg",
	"OlOjKgRVygySrFKXJNCIqjbLEaQuprbL2f1lDpMEMhyCklva1CFZFB9A7A0+NC+G66QOBaCeu/mpp60K",
	"XHqsmrFucbOIcbi1/6lg1Gk4b2EZuprzXc7uoiLWopHrlJyvgzCXKOJ0shIWSiA0IuwSjYuyrGOnVBp5",
	"/fvO9e/tdejf5RbaoIK1lGq+tUnVXGLjB9XQrc25WOWSrK3c0sby9mptq1r7d4EexpdFIqYLzYe6ocEs",
	"RTEKGckX9TxS7cqOS6jScl0rCvUgoAXNUBpxhbq+mn/MEJuh3AimaU4SgV21lSBfDtPdoPGCkBjBtN+2",
	"sPvLPVEZUeG6LKuYElWytxrV7ARb7rkNBYP/RZZCfVc0RD//BS8hQ3tKn35nxGWEprCImd6ybmQzAswQ",
	"AtmFQ+yZlUggC2dC6TYnikQzPL3zKAjFLU+7ubHUihnozvQ3z2oSrquFZWa3tOkwUTo7mUq5UHeK9qpI",
	"2SDK5TDtCK/XN9FdbdGh9Ld2uVlJC7J2iQ22KbIhVYU3TYDKmjaW1tiQMaTMtu+pW2zxZuUu0ixPiJ9r",
	"lCOQEmZUdceyYoaSpdhw5cihRCnMczjvq1dSdqZMMwWvpWS66kqNs1nbeffmbF4i9kjhpY/9rDErkCqQ",
	"/jCNaHvSet2rJnVXCxemq2062FxnJ155O3wz+AWlDxLdHXjuRHA7ZicOzDZap+jrSijl+g/v28UX7her",
	"HLoWrDqrjKXWrOwwT12t5TFR9DGN51rwriA0lHZrVEP90+iFsqAmKMrS5hFyQzzIxfwdplHcsvAzUScO",
	"ECqXWyush+x7Jkc006qX6tlVyzsYRqOhLDiIi0vcYqVlok7aTIyABKbwEq1rfnLw+vxMaW1+qryDwBoN",
	"lXb1ei5lab+D+dYTksZRwTW0FbjVUKH6n13MzwwzkOhw1WiUNOs6DB5nY1GISdqHkdWmLATDvfMsPYEq",
	"u6qXVhBmyjtYf6Oh1vUOchIVIes808hkm/rZRj83jlvc6ZypD5+pDzfueBwN6nc+jSaVYwv3HZCrT+ud",
	"0Ht1BUSLJIE5/qZP+TkbbroEbOyy52XrZQ8XTPNj8sl1o2fswXIxMRXXXdLOQ4Zl5ogykiMAgX3htoJ1",
	"JwY/Y+SssC7xaoXWzUVZ3KGy19vliEFuye62WD/mHMGcYYFwHsaIGo3GLKCxdFBaJMHO52BGijzmTCyC",
	"WPz3GqEv4o+EpGwm/pojyNucOg75lzd89GTOILXwZZeV6CpLbWz9ex1blWaCcN14ctD0zkkKwE/g/ACl",
	"EU4vz8EQHFeIJ5MVgBuEMeKDqR5HRRgiFKGo1kc1RBGgvAWlQpLpY054BXEs7l65RlJQpEb7FeK4MdRU",
	"FJoT0rSgvKfqsYc4OA2QZ5CCC4RSkMD8i7o5iZCEfKAuszEFOJXXfYhyjGtiUFgIBoGZXTAIJGzBINCf",
	"bFICXxE+yNBCulhz00WuwG2AVGBUBnJjKBiYmVRaO5ezMtlK8yVXUtFoH60YVZRiVNOJq7zhVVM9fjAa",
	"uj4llY7on8SZfJtL5QU8gJRekzxqkdKq1qjqe693xTneCLyBKSBpPOcMvLyCuJ6hVCp68gpaM3LXncB1",
	"jhmykdVHjF/AMw1UKbqrhUZc28UdcqzWXVRRlFLM8BWSsPk72nt1UE8gThlKYbrYv/V92dRcyPS74rV0",
	"jvVdYuIkKZhgf/1pfKkL3+2NXPiWp7h3c9e7CpoejYPmQ0Tehrw2+9+T75W2o78j731H/otED7eafs1J",
	"sviW3DRd8rpbL0/tqpsJEnCq9Si/wlwiYZQb7cFoA+DtFJAEM4FIZm2aSjehbc4QVxOiUeUUJZ2MaALj",
	"OGjVeacCHcE4jAvKUD5WA/NxaT8BJuam94f6oVdB/Ow4R5L1tTvz5R1qNdZb7+G8SuJVEq+SeJXEqyRe",
	"JXlQbnteJbl3leTVnaokC732YB7O8BXaw/TLEf7WdviLvxlhpzqAKxIXCaKOFXojtQhx6qRonhM2TsMc",
	"yTsj3VsMXF2f7WNcVzZE/OSzk5PRZxk6+Y/nf5lfPz9//uzZ5/96/9vxwf4pfv7X57RIvshfz/+x2qWL",
	"muBZhOmXMw6gWUpXjV7XZl3HlZyzsSLYY4zyA5jDBDFOcO33NESTuCT5zPSx3FXV3sC1W7vKXhE3IIJz",
	"mQ4j14XOcueAbq7z4mYQ/IsUeQrj/hSnOtwVxb2a/HbfJKdm6CA5V41GbbOuw1PC2TjrcQOsxOOYwUu9",
	"IB8Ksvca4EQ4ThAux0bgV5Jbx/p8egNAEQIzxjK6Mx7PiotRRMIvKB+FJBnn47Qg0YX6lzdvcnZDrlpC",
	"X+M45utsWHuD0ke10OstbqGIJfzL+tuspLOwXOdnn0dnQ61C8T9/XnWB226z2y+xu+6uXza0pnpT1mAi",
	"babnAhvyYwb/LJDNXGiRZTGWG8zJUEarWX2861n5nYr8rJbbotSu6fCjaTR9GK68Lt2oIaFDeIAS9948",
	"2H8/RCmHOAIh7zEVWqjgkEf//Q6EMea6GF+qK5Tj6by6YgqdK3iqwLMMJaWLiv5pfFNkQQc7Klvczrui",
	"6YFw7/4V3a4VZgG0c8VtfCeqThN1b4kFYZ6yAZ0VLCLXaU+oZ5AC02XFwJeyu7o6tcbTl6e6qOOQx25D",
	"/4z30ygjOG0xrZGqdW2NkKSpMLtX9Fmif8ZnevxyUtVCMzG7uENDrLfr65mhF0r7ZuxqkaxdHRr0xzEA",
	"wxBlTOJFIgOTlKoh3qhrYtcIF0g4BCiPQW7XpoSBOWKlLqBGeU8iPJ13DpPwJhhF2kmEkSxTHY7UWtc8",
	"G5QqiKkZze5b+oRUyFc4T1DZQDXf/5pxU8vZHMk6p69JpaHypxDonGF0hQAE0qdCLlGL70kTtt6OHWpE",
	"aZurIZWhLjWj0nejD9IOCRML/aaUJFRA+e7Ili4U5KIhWbwalnuMIcRgEGiKCgaBIYtgEOgF138Kvxm1",
	"NE5XmkFgps7/doC/hLdNCWAtU9XijWJPqdm5xx6xJ1UbYCV3H400x1iojk9Hm0WUXFm3tgnr3VxZpZ1g",
	"FfpsW9yd4FakqYmsDePUEKGhzJ2gNyf6AV2a9vOc5G9IypAUwzU9lkSuvCrcmgjFmg1jdIVigPgogLeW",
	"h6my/QWiYEauhZCTLcqDVen+L/eWbnMuk7GcgylGcWQ1xilDeZYjhiItJH8/Pj442z88/HiomXLbF/jC",
	"RkDcbyDAuwE5ewnus3P56/y5dZickZQizjhILrznGAGHv74Z/vLL1kSeBBtInTACSAGUuWeGJveM5GUj",
	"LaA/fnjz6fBw/8Px2aeDvd3jfT6L3ebRMQihdOHl6BHn2CQH5we7x29+L4+0GRHbfwR2a2fd6m4wRyzn",
	"1h+cMpSDQqQwOf9t//ic9yQXDGKJnJhvR1beL2obns8YZlk81+p8hChnSCCcwVSeJ2Amv14FrPb9a8xm",
	"pGAApnPVlZZ3mqqHpESNpk8f/uvDxz8+nB3u//en/aNjvdRSETSdhCtkTuS1CIgKARFMQZF+Sfm2V4MK",
	"//4BSBCbkWjAEWlmmkE2G4FjzgsUzPp2AECNBj55TGmBwAVi15zZsBIUjiJpno0sAVoSKRc39SUPBkFt",
	"fkuIvuZoO4GnH0M/wcBGvrInHgWDcFCFBP9hE33FLhJHUnz+TdmxJ8qtdFtiVdxGnLKCK9dW1dEaq8VB",
	"KjMotQ255UyVtP81g2mEoneYsv2U5fOmPJSXZk6LrrgYcoSCHMWQ4StUp8xPh++aDhvV/aNXpgv0iRP0",
	"ZurRFUIH0VdMhW+x8VoXoYR8d4rhg/q99Xoup9cVZjhZd5ihVP/t4JoNoeLWEYlbnRGJtYvKCmHU0HS6",
	"kk4qcXRmiFlOo16qp1Etrymqjcq3DCV8gzp01SKnpOV+WtYByLT1xVAiQ6Kl4RNjKq4wKqKgoCjX27eb",
	"m1hxRo5cDZgyLjrkN3PEijxFkbg6QZzjaEcG5ZgglV8qYhzVx1GkhC6H6s8C5fPycB+ci0EiKdjl36OT",
	"YjJ5EYqBxJ/oXNzOyGlJOPS9jFxXOgCEzVB+jSkSTeemgZwvhyXLEUWpyYunOR4t0VrxTlN8bkqKVNzn",
	"GxSRFH2cBjufmxcZ3V4JTf5847pGb4ZolfeYMU4wawsq+4qTIgFpkVygvFwwybOFEJ/BKyTNTL2KHRQj",
	"FvhcfPCcC2v0ZwFjLRTqHynHE8uEhXjPCKVY3MxxEBJugZd0W97bPaMIgfMUfWXnz/XpgsK9lC3SYwVG",
	"VzAN9eKdk+mUInbOazSQXDkqbDVNJGC4gnHBNQNZwmn5GckBZhSc88U6f86VmXO5wc5HdrwtTtmL7XLj",
	"cP3nUvhVVEQ7B9y9IkJiVgSkyF9Q2U8DgEss8y1T+v5wVaANaSNhRsALitLQnMNKHCrjs+YdBONYjdMk",
	"glEXc3h5Mwgkrls88ETd7TjTEhjnqhkjDMYt0da8qkGdJdqW+9grp67yT0rSD+pcocEKTDMXc1AlEjy7",
	"RN9lOPkB/94BZOHsY4Zy4wtZT4Ep3anKe+jx6KdFPJ9kdmApjKKAC9aEXCHxRxZDIedUQUgycdvDl+y0",
	"U8ET1+EzPnabm+AC4Pi2Ett2kZ+XWYq6TkBUPvSZ01Gp4XXbwKc8saS7bEFgf0VyaM8XkpcHe1oKRVic",
	"ZUZrDPhf2hGGr4ya2Nu0Y1JFynDcOiNMzWQGAI0uR+B8KzqvzGsrqrpEnJxEPz8/OaE/aSfu//t/Tn9+",
	"vp5JCeWF7mn0dt4ouqZT2re3vGLE9MxaZKUyVsqMwmiVdlyXVpqt5jtge6Kb++ZKmblztkprKmyt6kAi",
	"0ccdPLh0tRuJOyhdrx9R2MFmUt9q9uKjDpZzplccxfvSLyP3/v1mENAYulFy9G7XeKxqL0XwBqYpYdZE",
	"AElDVPXZNZmYqu6qEbq6oxMjPgVz9yf+Nhd/Mez0xJHV/eIBFIXdLhyg5HxOz/5F2LkPP/+XK/v5q7Xg",
	"HVwadMuWdeXnfcDOvprs676+fOvHV/WbFb6p+M5Xa8e39hRfFtIk+zGcgLceihOwKwDnQfgATx6FD7CD",
	"H27WBVhchapDMSuYbyNewZO78Qq+1TTW/cJOD6dgS+Zv0id461Y+wRpocSMfx/KIr/R/Uqeb9+UsvN3h",
	"LNzTl1ZNsNWV1lq1urdr3WfWanoLl9nGKLfxmFWDuR1m7bV1LOpiF9oFA3Q51equt/CptYdYzaV2fa6w",
	"JIHipvWOHWHX6/zqIu2ml+tyVN10cl19mRo+rv3ore71ugypOZxeW/fjWh1WNQFV3VWX2KFLOrBWw/fX",
	"7b/68gH7rzaC3FscX3q/Q2UlS61E0XDTU5hA0J2rGVP56oLSsq0Ld9s0HOilEq5YroQktW3JdzOlRVI/",
	"UoBhgsZSY8Hp5ThCCRlvT7ZfTrYmW1vbk8lkUrOoKgdk4+ffJ4MXN8u8c7u6P0vNkaXhwbLYi2dVvV8t",
	"5Jk4Nih1rEphqWpZxTd1LataV3sSyfEGJ45b8nyXd6uiUfl2Byc7JrlGryvWxVmZZX5bnYJZ/irzLfPf",
	"nQ9sqgYqYe6CychWdzcbk7ZXTsf81PNRBR2UVLbQiX4XTEk1u7s5lQmHldFifhujRZV0GS1lE4oYl14L",
	"T/MOlfg6Uu3NgZ5KhbwAL7LV3aHFJGSWWDE/NVJUQYfhULZQmZwXzEi2ursZKSj0jMxPPSNV0CFydYtV",
	"2aB+6s1KQG1KrPTTqqzB/awKN/E4OOD84/QPhL5U3rQKjoo0El5ZzeWI4FyLQr5+xgszy0lCmC2eDUVo",
	"9fw9UaMeF4jKv/5AUar/Pp4Vufrz1xzLP44gK3L1p4TpdMXc9/MzMj3jIFm81i4rOW5Z2sGmas3E9q7i",
	"8J8wLWDuRqJortHIiaYDjYYMNRrLgX9FF7n68z3Mw1kwCHazHMfiNy/9Z5Ei8R8xwG5xWQhHpCOUMSSc",
	"cwbBx5AR+dcHcqUL91Ao/1wN2xIbFX5Z55ZdMs00UMh4J+IVjsmeW1xXX1sT2LOiHPRTFOVBHKdgPAVJ",
	"ETOcWU/kCr9qpatDBiTVrfiyGofiTEJwxshZVci311sHnM4W3Zkq2rrU0Ph7i6awLB65pO6DSN5ubZis",
	"KRgdDdpx6VBCXixCZtmnhs33bVrKsuiUPKEHPs1+WgtC69pNV4t2lLo0oFeLcGp10tEYx+QdpGyZBxZF",
	"wJJtGBqrX+hRqBHoIbxly8z5JlA2pQxB40Sjb7OrbmlLHhnLD/Op8o9aUr1RUQr3WlVnnol629XUDqOP",
	"GuO+LDDmvS6qG/hWefUZ2QfjMnXc8u4mI8oXyjrfkbnfHmU2VhqSrO3knVe1PV0NxTmI8uKxDyJ+Ov3L",
	"dRqxfdPnFGIx0QlwDcWpX4bcxO+uHB66AY0hbXUO6VhkE+WwhqdjBAyWdwetund0vj6r66VrRIuzKsqX",
	"n0k/X4yV5qvdOMoL5Pq18cKEUbThnyFX9HRFDqYeWi45mCkoOZgqanCwsly+FLJXPgnSdmcshUNrDJOQ",
	"HHuvd83jIfqkPSmo8Ewyd/xcqIu7YxWnD7IizwgVt92LPG0sd+cc0XkaduogOs5YeXOgqAIgoEQyCMxM",
	"0kIYsgLGlWYyWIZL5POBCCbhMxVhM3IS0lFMp++QAbeVQ9MI5yhk4smUKcmVi5Z8+MYJ16hF9pYO9Azm",
	"l21BG7KuMuDCcFObJvU6N12DxDFOWOSYzY+EmqGOsCkOdwtpCUr1I9gJXu8evX1TfnfGmDhKvUAwR3mz",
	"9f7u4f5hvfmNSGI5JTLZQ8qgzHOMEhG8GwhHotFRkWUkZ//5IqKjUByaFnmshqA7Y+kjJGoaQlJ6Ir0h",
	"KctJDA5imCLw7M3BcwDjmFwLH0DJg1SwpQjakc87yq4lU8oRV/PiuYqFg0BnYpPR1EfKleXZ3msIj56D",
	"hG8xscwoT+jHqaq3wI5IODKgK0cnqUSOcxQjSNEwJQxRWTWMcYhSioZivLHw7mLiKN41x8P9o2Owe/A2",
	"GATau2cn2B79fTQRHrIZSmGGg53gxWgy2lKxFmKpx1KKCiGquIIiQ6LDV95GwU7wG2K7cWy90Sw9WSxf",
	"mZ3PfeKNGJEBVZLbc2IIdgIRY6ifn9/R8UsDSUzCZ3LRudjNzaBPHGbl65qZmUhDWurVvEE1dtI4pg6s",
	"4CzdwjDEGH3FIbnMYTbDIYzjObgUhCZMnlTFMSqeKUJmpKMcUqGG2sOVSoc3GfGmQ+KoeDEL5shW7ikR",
	"jlQuNKoAVRuNNSXBjbRGRCJR8xR40GGjZSzpNaQApSLwQeYisMLjxBVViTmZjEEYRHz87Yn2yx2Bj5V4",
	"UNMKUzCReVpVqoBaeFzp7muHx7nQIWMxb0lUliCSa1YLo5WCRVb9h5Iv4mKOQ1yLuIU5Mgtf6dcIpx2M",
	"RiMhnRy3fYY4GvkeQpIkcEgR36O8UAcyMpKpNDoKEDMZG4pmUK+C659HHz8cQDbjfXJEOQmooN8WCB0d",
	"jDu2vP0k8RUygYbDUqng325bTgnlctT9pgUjUiGlpRluHEmlXaI5BV8xccFBTYj0+e6HvfMR2BVJv+Rz",
	"dVK7LScrVnrnJP0JnH9B83MwBB/TeK7o1fi/K1d2cTMvvirOVCx/+C9orsf4D7HotxkJcKbMiKQeMezf",
	"FsMWEXCWEnbWD8i/9YKy35gVcNv2N+/7q0D+XVOF3DdrowoxHN8hfVYUgpSkw/O0iGMtGZQ9XGJMhhRz",
	"GV8dvz/dwPLR037f+FvfSUREf2G1OaxCVsvMpoW6RINVqMuSGApUmM5FkkAVbU+Li2EJPbpCqc7FbiJR",
	"KLeWCZLiVM5I9AdZjq9wjFRiIkG4XNCUw4mJYgoEv28VjJTtGnhcszPxzjenwkYT6YiErrg9mWhVXhmS",
	"sEzWNv4XlUEE5YBd98Qm1YawEtoTW9jnPhjpzBVS+HNt9+UaYaqkq2uBy46OF8Dg9ArGWMGytXlYIixd",
	"67KcXOGIM6c8F35gBZfJTJvoYY5EhBiUB28vJy82D6rWETloJMfftHedoG4rnRGKymxFN4Pg1QaXeFc5",
	"wat8XiQU2zKq2M7C+LGs5s+nVbv48ynfOdIffR7sBJzMG2eYnJaxTnglzri4SioOiHTWHGOmnXJzv8Ay",
	"IYAKpRavm3Koajbd+LvtEXbTZeItsu8E+1Bc0hVHaB87yPuOJW0Pby56c9Gbi95c9OaiNxe9uejNRW8u",
	"enPRm4veXPTmYpu5CFNQTxCzdpNx/F062dzI6+kYydjYqvkoAveQZUHO79x+dAyX6S/f0hAtrRqKQpJK",
	"tfcaYmbyi5jZ85oLLrJTGHOSHIAijTlBCnOlFBoqSNGkruGErAZtYbUMJ4gU7EhCUGG1xq1w4ggaSHCK",
	"kyIRlQ0jpcmTX7q9Dqr+VXzryaX3XPPH4JovJWFsDNYqvRkFRngNS3j+fo/waCSqPWB7PJdbXG3YpyBz",
	"JDvn4qWep92w2KUFzaDPweNjERt3qdva6FDhcn1pVOS/9hzac+h74NCPnef9hti6GV4GWThrsjyRIPqx",
	"Mj1Biq+JzFjUsdhDMfmf9bqrZAKCWQo6KRH/Rj7Mo/A9FZ9IBWgqr3Tw+SQF4Dv/B4CTgGQnwQ44CWAU",
	"nQQDXSpOf0T52AxhVcvDS1H/nyKm7CTgVTcn6Skn3q0qSEeItQchy3hdkTZjRQBNkOtYDu2C86UN3nYN",
	"vIJmKI00RKsDQuVAKHJBwMnBBuJFFYhDRIsE9YZBJitfCEb5wRubFE18wKIM47Xk6/XQgJuqd7TYqA9B",
	"lIvNIiQBLBPnNeJcvGT3kv1ebK9f7hEe+RxKLRuueoKsSNHXDIW8xGQSern1aoN6SIIiDAHnMfJRG3l3",
	"f94hDdXTIeri8iloTjLiaO3KU+GwFkWCMfTUdadFCpMMX6lukymxg+rUAUXzfFgLaSGfTwJLVRKqUUUf",
	"OTFxZbzWSHQR5KYUizBBRmbfOBSpJUC1knKuH8gye1UV3O2VwOV4PXq3uxhSpWl2AhpDygs+n4gUvSfB",
	"qQ1fTedSO00F/3LAYGOr1UCy199eNAUTp2hZrA6VTY38SAXWPCOy7ZbW1Aarz9tJRapHLZe8+uj2i+0t",
	"3u6mphyuoIE9QBWweZqjhJ2g0zVqWreASCWK83qo10PvR++7QOKNYYDZg9NTFaQ6lyLJH62y+vQ0VCXX",
	"zZo4xOYdXpyP9RnJQudr+mtOksd/f+4dub0jt3fk9o7c3pHbO3J7R27vyO0dub0j91ocub0Dtz+G8I4u",
	"63Uop4DBLyjVj2Nu0i42p+pdlvF7xcX3TGNvGXvL2FvG3jL2lrG3jL1l7C1jbxl7y/jHtYwt1xRvG3vb",
	"2NvGa7GNy13V58mBWxjJ/bIu+3TL3rj0xqU3Lr1x6Y1Lb1x649Ibl9649Neu3rT0ebP6m3KGmHUSp7ZE",
	"yz0NtyVSK3+0GrYadD6/srcPvX3o7UNvH3r70NuH3j709qG3D7196O1Dbx/eg324KLPyajbi+LtKctFp",
	"LR7INpsyFF3uriYVhzc5vcnpTU5vcnqT05uc3uT0Jqc3Ob3J6U1Ob3J6k3P9JqeVAXG9xub4u/5gjwdg",
	"H5PB6RhJz9Tbrt529bart1297eptV2+7etvV267edvW2q7ddve26pqxB0mS1zK3+NusgyAhd8LrDD2iF",
	"3t0bEUItkQph84UAR6J/wQy4RiEz74uXoKx3E5bMvi9hWSL7/rpz3ffJcu/T2/vY+Y3Ezlfy2uts8Tqv",
	"feVw8MFkuHeC7LPaP5Cs9jAFJB1GKIHmScC7PEEef5eNe78H/2OeJzuGMmvjn6df0/P0/l16L7Q3IgEf",
	"zIP0/iX6VUTcorfnvYhaXUTd/VP4fc03//q9lwMblwNP89n7ZY/zFj507znsLTns+t/d340idWvJjQQN",
	"1oovyctDxLF1cOh6er9ykOh6eP9QvBKvwJrmJOkLWNfz8k3YnuIr87d5X94LTS80N2s83f85pn9K/jE9",
	"Jb+sQtL+ePzHXH7E6yUb1UsWKCPNl80bgr/fg+YKubLGevVdVWuM6TfPE+J4Cn17sv1ysjXZmoj/db2I",
	"3n5R2vcl84YK9DbJSO7SycEMppHwcqlhBIseR+I7Fehkv99FtxJVQ4OVIZ//0DnbsvtBXFxiBThKLlAU",
	"oWgUZqO0INHFKCTJLW6G7/FF9r6HCnf6CLu/mPZq2qO/mB4Akj+e19qds/SPtT+Gx9qXvtVWCR7HIUxh",
	"Ps9JHJOCcXnSXjn+ztWAm642DCVZDBmivRq5B9Rb6c+CMEg7K90DzFCcTBFkRd6Aw65yd85JjNpmUalz",
	"d+dLjkPEsFCZW6uqnXu99bcbxx3P/PkQLR+i5UO0fIiWD9HyIVo+RMuHaPkQLR+i5Z/T8ycvPkir4/k6",
	"5RzYfOugNMi6jWjTbsX3Djb3bru3D7196O1Dbx96+9Dbh94+9Pahtw+9fejtQ28fevuw2z5sf+vgtjZi",
	"v/cONmch+qcOvK3pbU1va3pb09ua3tb0tqa3Nb2t6W1Nb2t6L3CkHx94giHuLVav47mF9dm7tScXuvNk",
	"7ZVhXP7lBZ/g6tiOtfAprjzjv1tYO/NN/v1+QPFprlyZlfuIp8HCY9YfNbPyHSnzGq2dsZwVwvZ5qjwz",
	"vwdm/tQyVa3MILuyVf3oTHL9aabezGB6aREmw+JgacU8U7y3O7dUujWiCYzjzvxSb9MwR2J7YPoFUPwN",
	"CZPCAAfzcIavEF0dwCznVMUwomM12B6mX47wN+QGe2sy+Q1XYd6uwnyEmIU9cYOBc7Q6hAnk9kgK09AJ",
	"khqIl8kv0bc6bYXKnCUgLcF9UQV3D1NxDmzR4J0Diqn6qsjDIbJwNOF8WYVzP10OzK4MYxVIn1ZysV7q",
	"TWt6MYNdr+14bWfDpusv9wPK008y9nJ7+wGh1koyci3q0dcQoQhAQ5lAJFJ4ounRVteC21OkeSV4nW8w",
	"OTQMlbhLKrNCeZ1o5VXnKbuAB5DSa5Kr1GKhUKPfMjvRmVI0Kzm/Gkon7/xK6JhGc/kXKfIUxrU2lhrq",
	"yktWm430XdDOSyUwzVl2TEV9aduZC64db7dLAVfmeaukf7MWY6u2GMugWunz3bjeKnX+rjxxLyavJi8U",
	"lnpnVtObV66XxOW95Flb/mDsjnKtLQ+Iz7fmtdZN6zOujGvuO/lHknCte5pPKuXaw9CIHaj9kZRjR964",
	"FZXjpV08xpaOw8dzP28qZfGe1dT7e/wI/h5rMifepphhTtd7r3dBpkhI0XpTR1YulUrjlj/2bE1cq70w",
	"v0Sqmfy70sptCxwiOk9DEJJ0ii+LHEUVmBYAA0NWwNgFSy6GLY9wXcbB/6AcT+fGY3SJ77qQsJxq3djA",
	"y+rVTdKvrKWlAwNaCLVlWsTxXCrE2z1GwCbhR5YTofaQKbhAOL20lWuv03qd9uk6ETW2FIdGb6tFjkSb",
	"VJJ304bOq0LMQlKkDOVqt/q0wxs8W2V1Cmp9576vJomuxNxkyuKMxDjEJuuus66ZaKlv0/F3UTF39nF+",
	"s8/HOqN6l2hd9YleraP9+rA1gp1d2FXcOcte0ctLtW+dqWpA3aWdULY0agOSCw/qKOr8hqvF+Dsv1Q1n",
	"CMZs9k39isklV6XF3yRDKcyw+mVNtCvb84Fu5pM9+wBrH2DtA6x9gLUPsPYB1j7A2gdY+wDruw+w1lqq",
	"j6/2ubwebVSzIeLWVM/GGOs+qGg3Q1tNuHb7zed29uagNwe9OejNQW8OenPQm4PeHPTmoDcHvTnozcEN",
	"moPtmZ1vZxJW8zp3J7c6MA6Hjze3s09IZVzRfT4q70p2X4kIN+1JZtP8MtmoNh1poeGsRANoYK0QgBnf",
	"u8ZT5mmmzWomdewh6gaLDjgfi/S6Sy1aoaI7EYS1Y3yaKy8jfLLa22a5WpGfdeW4eow87c7yUmlKurO0",
	"VJMeaamOECsBeXgZnrbdGZ7KxXwYCZ5eOBM89YXyB83v1Euqt6V30rj1Qt4L+c0agg/AwHqSuZ2eYIak",
	"VTWo9vxIT1mB6pfTqClUVcQtjaHOrXPlSq0zGaUwJW1xxWr1WkdfkPinzOzTzPqzGLBqzp/1JeRpStgH",
	"Jt0tm/2OMvAsC4bPv+P1inuS5g8qhY4T0ieZQecJJp9ZSe+4GQSV0M2OmM0FwZpiMnyaUjUp8jjYCWaM",
	"ZXRnPFYSfpQWJLoYhSQJbk5v/l8AAAD//8DbqyoU5QEA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
