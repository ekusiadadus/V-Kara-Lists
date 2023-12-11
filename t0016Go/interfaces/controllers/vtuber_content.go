package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharin-sushi/0016go_next_relation/domain"
	"github.com/sharin-sushi/0016go_next_relation/interfaces/controllers/common"
	"github.com/sharin-sushi/0016go_next_relation/interfaces/database"
	"github.com/sharin-sushi/0016go_next_relation/useCase"
)

type VtuberContentController struct {
	VtuberContentInteractor useCase.VtuberContentInteractor
	FavoriteInteractor      useCase.FavoriteInteractor
	UserInteractor          useCase.UserInteractor
}

func NewVtuberContentController(sqlHandler database.SqlHandler) *VtuberContentController {
	return &VtuberContentController{
		VtuberContentInteractor: useCase.VtuberContentInteractor{
			VtuberContentRepository: &database.VtuberContentRepository{
				SqlHandler: sqlHandler,
			},
		},
		UserInteractor: useCase.UserInteractor{
			UserRepository: &database.UserRepository{
				SqlHandler: sqlHandler,
			},
		},
		FavoriteInteractor: useCase.FavoriteInteractor{
			FavoriteRepository: &database.FavoriteRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *VtuberContentController) ReturnTopPageData(c *gin.Context) {
	var errs []error
	allVts, err := controller.VtuberContentInteractor.GetAllVtubers()
	if err != nil {
		errs = append(errs, err)
	}
	VtsMosWitFav, err := controller.FavoriteInteractor.GetVtubersMoviesWithFavCnts() //エラー発生
	if err != nil {
		fmt.Print("err:", err)
		errs = append(errs, err)
	}
	VtsMosKasWithFav, err := controller.FavoriteInteractor.GetVtubersMoviesKaraokesWithFavCnts()
	if err != nil {
		errs = append(errs, err)
	}

	applicantListenerId, err := common.TakeListenerIdFromJWT(c) //非ログイン時でも処理は続ける
	if err != nil || applicantListenerId == 0 {
		errs = append(errs, err)
		c.JSON(http.StatusOK, gin.H{
			"vtubers":                 allVts,
			"vtubers_movies":          VtsMosWitFav,
			"vtubers_movies_karaokes": VtsMosKasWithFav,
			// "error":   errs,
			"message": "dont you Loged in ?",
		})
		return
	}

	myFav, err := controller.FavoriteInteractor.FindFavsOfUser(applicantListenerId)
	TransmitMovies := common.AddIsFavToMovieWithFav(VtsMosWitFav, myFav)
	TransmitKaraokes := common.AddIsFavToKaraokeWithFav(VtsMosKasWithFav, myFav)

	c.JSON(http.StatusOK, gin.H{
		"vtubers":                 allVts,
		"vtubers_movies":          TransmitMovies,
		"vtubers_movies_karaokes": TransmitKaraokes,
		"error":                   errs,
	})
	return
}

func (controller *VtuberContentController) GetAllJoinVtubersMoviesKaraokes(c *gin.Context) {
	allVsMsKs, err := controller.VtuberContentInteractor.GetEssentialJoinVtubersMoviesKaraokes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"resultStsのerror": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"vtubers_and_movies_and_karaokes": allVsMsKs,
	})
	return
}
func (controller *VtuberContentController) CreateVtuber(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var vtuber domain.Vtuber
	if err := c.ShouldBind(&vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	vtuber.VtuberInputterId = &applicantListenerId

	if err := controller.VtuberContentInteractor.CreateVtuber(vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invailed Registered the New Vtuber",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Registered the New Vtuber",
	})
	return
}

func (controller *VtuberContentController) CreateMovie(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var movie domain.Movie
	if err := c.ShouldBind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	fmt.Printf("1-1:*値よな？ %v \n", *movie.VtuberId) //4
	movie.MovieInputterId = &applicantListenerId

	if err := controller.VtuberContentInteractor.CreateMovie(movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invailed Registered the New Movie",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Registered the New Movie",
	})
	return
}
func (controller *VtuberContentController) CreateKaraokeSing(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var karaoke domain.Karaoke
	if err := c.ShouldBind(&karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	karaoke.KaraokeInputterId = &applicantListenerId
	if err := controller.VtuberContentInteractor.CreateKaraokeSing(karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invailed Registered the New KaraokeSing",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Registered the New KaraokeSing",
	})
	return
}
func (controller *VtuberContentController) EditVtuber(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var vtuber domain.Vtuber
	if err := c.ShouldBind(&vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	vtuber.VtuberInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyVtuber(int(applicantListenerId), vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.UpdateVtuber(vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Update",
	})
	return
}
func (controller *VtuberContentController) EditMovie(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var Movie domain.Movie
	if err := c.ShouldBind(&Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	Movie.MovieInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyMovie(int(applicantListenerId), Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.UpdateMovie(Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Update",
	})
	return
}
func (controller *VtuberContentController) EditKaraokeSing(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var Karaoke domain.Karaoke
	if err := c.ShouldBind(&Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	Karaoke.KaraokeInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyKaraoke(int(applicantListenerId), Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.UpdateKaraokeSing(Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Update",
	})
	return
}
func (controller *VtuberContentController) DeleteVtuber(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var Vtuber domain.Vtuber
	if err := c.ShouldBind(&Vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	Vtuber.VtuberInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyVtuber(int(applicantListenerId), Vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.DeleteVtuber(Vtuber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Delete",
	})
	return
}

func (controller *VtuberContentController) DeleteMovie(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var Movie domain.Movie
	if err := c.ShouldBind(&Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	Movie.MovieInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyMovie(int(applicantListenerId), Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.DeleteMovie(Movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Delete",
	})
	return
}

func (controller *VtuberContentController) DeleteKaraokeSing(c *gin.Context) {
	applicantListenerId, err := common.TakeListenerIdFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching listener info",
		})
		return
	}
	var Karaoke domain.Karaoke
	if err := c.ShouldBind(&Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}
	Karaoke.KaraokeInputterId = &applicantListenerId
	if isAuth, err := controller.VtuberContentInteractor.VerifyUserModifyKaraoke(int(applicantListenerId), Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Auth Check is failed.(we could not Verify)",
		})
		return
	} else if isAuth == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only The Inputter can modify each data",
		})
		return
	}
	if err := controller.VtuberContentInteractor.DeleteKaraokeSing(Karaoke); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Inputter can modify each data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Delete",
	})
	return
}

func (controller *VtuberContentController) ReadAllVtuverMovieKaraoke(c *gin.Context) {
	var errs []error
	allVts, err := controller.VtuberContentInteractor.GetAllVtubers()
	if err != nil {
		errs = append(errs, err)
	}
	allMos, err := controller.VtuberContentInteractor.GetAllMovies()
	if err != nil {
		errs = append(errs, err)
	}
	allKas, err := controller.VtuberContentInteractor.GetAllKaraokes()
	if err != nil {
		errs = append(errs, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"vtubers":                         allVts,
		"vtubers_and_movies":              allMos,
		"vtubers_and_movies_and_karaokes": allKas,
		"error":                           errs,
	})
}

func (controller *VtuberContentController) Enigma(c *gin.Context) {
	var errs []error
	allVts, err := controller.VtuberContentInteractor.GetAllVtubers()
	if err != nil {
		errs = append(errs, err)
	}
	allMos, err := controller.VtuberContentInteractor.GetAllMovies()
	if err != nil {
		errs = append(errs, err)
	}
	allKas, err := controller.VtuberContentInteractor.GetAllKaraokes()
	if err != nil {
		errs = append(errs, err)
	}
	if err != nil {
		errs = append(errs, err)
	}
	allVtsMos, err := controller.VtuberContentInteractor.GetAllVtubersMovies()
	if err != nil {
		errs = append(errs, err)

	}

	allVtsMosKas, err := controller.VtuberContentInteractor.GetEssentialJoinVtubersMoviesKaraokes()
	if err != nil {
		errs = append(errs, err)
	}
	all, err := controller.VtuberContentInteractor.GetAllVtubersMoviesKaraokes()
	if err != nil {
		errs = append(errs, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"vtubers":                         allVts,
		"movies":                          allMos,
		"karaokes":                        allKas,
		"vtubers_and_movies":              allVtsMos,
		"vtubers_and_movies_and_karaokes": allVtsMosKas,
		"all":                             all,
		"error":                           errs,
	})
}
